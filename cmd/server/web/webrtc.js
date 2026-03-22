(() => {
  const START_LABELS = {
    publish: '开始推流',
    play: '开始播放',
  };

  const RETRY_LABELS = {
    publish: '重新推流',
    play: '重新播放',
  };

  const ACTIVE_LABELS = {
    publish: '推流中',
    play: '播放中',
  };

  class WebRTCSession {
    constructor(kind, options) {
      this.kind = kind;
      this.options = options;
      this.pc = null;
      this.localStream = null;
      this.remoteStream = null;
      this.manualStop = false;
      this.state = 'idle';
      this.updateButtons();
    }

    async start() {
      if (this.state === 'preparing' || this.state === 'connecting') {
        return;
      }
      this.manualStop = false;
      this.cleanup();
      window.GoLiveApp.clearLog(this.options.logEl);
      const room = window.GoLiveApp.normalizeRoom(this.options.roomInput.value);
      this.options.roomInput.value = room;
      window.localStorage.setItem('go-live.room', room);
      const token = this.options.tokenInput ? this.options.tokenInput.value.trim() : '';

      try {
        if (this.kind === 'publish') {
          await this.startPublish(room, token);
        } else {
          await this.startPlay(room, token);
        }
      } catch (error) {
        const message = error && error.message ? error.message : '未知错误';
        this.log(`${this.kind === 'publish' ? '推流' : '播放'}失败：${message}`);
        this.cleanup();
        this.setState('error', message, 'error');
      }
    }

    stop(message = '已停止') {
      this.manualStop = true;
      this.cleanup();
      this.setState('stopped', message, 'info');
    }

    async startPublish(room, token) {
      this.setState('preparing', '正在读取运行时配置…', 'info');
      const bootstrap = await window.GoLiveApp.getBootstrap();
      this.log(`已加载配置：${window.GoLiveApp.describeBootstrap(bootstrap)}`);

      if (!navigator.mediaDevices || !navigator.mediaDevices.getUserMedia) {
        throw new Error('当前浏览器不支持媒体采集');
      }

      this.setState('preparing', '正在请求摄像头和麦克风权限…', 'info');
      const stream = await navigator.mediaDevices.getUserMedia({ video: true, audio: true });
      this.localStream = stream;
      this.options.mediaEl.srcObject = stream;

      const pc = this.createPeerConnection(bootstrap, room);
      for (const track of stream.getTracks()) {
        pc.addTrack(track, stream);
      }

      this.setState('connecting', '正在与服务器协商推流…', 'info');
      const offer = await pc.createOffer();
      await pc.setLocalDescription(offer);
      await this.waitForIceGathering(pc);

      const answer = await window.GoLiveApp.fetchText(`/api/whip/publish/${encodeURIComponent(room)}`, {
        method: 'POST',
        headers: window.GoLiveApp.createHeaders(token, { 'Content-Type': 'application/sdp' }),
        body: pc.localDescription.sdp,
      });

      await pc.setRemoteDescription({ type: 'answer', sdp: answer });
      this.setState('connected', '推流已建立', 'success');
    }

    async startPlay(room, token) {
      this.setState('preparing', '正在读取运行时配置…', 'info');
      const bootstrap = await window.GoLiveApp.getBootstrap();
      this.log(`已加载配置：${window.GoLiveApp.describeBootstrap(bootstrap)}`);

      const pc = this.createPeerConnection(bootstrap, room);
      pc.ontrack = (event) => {
        const [stream] = event.streams;
        if (stream) {
          this.remoteStream = stream;
          this.options.mediaEl.srcObject = stream;
        } else {
          if (!this.remoteStream) {
            this.remoteStream = new MediaStream();
            this.options.mediaEl.srcObject = this.remoteStream;
          }
          this.remoteStream.addTrack(event.track);
        }
        this.options.mediaEl.play().catch(() => {
          this.log('浏览器阻止了自动播放，请点击视频控件继续播放');
        });
      };
      pc.addTransceiver('video', { direction: 'recvonly' });
      pc.addTransceiver('audio', { direction: 'recvonly' });

      this.setState('connecting', '正在与服务器协商播放…', 'info');
      const offer = await pc.createOffer();
      await pc.setLocalDescription(offer);
      await this.waitForIceGathering(pc);

      const answer = await window.GoLiveApp.fetchText(`/api/whep/play/${encodeURIComponent(room)}`, {
        method: 'POST',
        headers: window.GoLiveApp.createHeaders(token, { 'Content-Type': 'application/sdp' }),
        body: pc.localDescription.sdp,
      });

      await pc.setRemoteDescription({ type: 'answer', sdp: answer });
      this.setState('connected', '播放已建立', 'success');
    }

    createPeerConnection(bootstrap, room) {
      const pc = new RTCPeerConnection({
        iceServers: bootstrap.iceServers || [],
      });
      this.pc = pc;

      pc.onconnectionstatechange = () => {
        this.log(`房间 ${room} 连接状态：${pc.connectionState}`);
        if (this.manualStop) {
          return;
        }
        if (pc.connectionState === 'connected') {
          this.setState('connected', this.kind === 'publish' ? '推流已建立' : '播放已建立', 'success');
          return;
        }
        if (pc.connectionState === 'failed' || pc.connectionState === 'disconnected') {
          this.cleanup(false);
          this.setState('error', '连接已断开，请重试', 'error');
        }
      };

      pc.oniceconnectionstatechange = () => {
        this.log(`ICE 状态：${pc.iceConnectionState}`);
      };

      return pc;
    }

    waitForIceGathering(pc, timeout = 3000) {
      return new Promise((resolve) => {
        if (pc.iceGatheringState === 'complete') {
          resolve();
          return;
        }
        const onChange = () => {
          if (pc.iceGatheringState === 'complete') {
            pc.removeEventListener('icegatheringstatechange', onChange);
            resolve();
          }
        };
        pc.addEventListener('icegatheringstatechange', onChange);
        window.setTimeout(() => {
          pc.removeEventListener('icegatheringstatechange', onChange);
          resolve();
        }, timeout);
      });
    }

    cleanup(updateState = false) {
      if (this.pc) {
        this.pc.onconnectionstatechange = null;
        this.pc.oniceconnectionstatechange = null;
        this.pc.ontrack = null;
        if (this.pc.connectionState !== 'closed') {
          this.pc.close();
        }
        this.pc = null;
      }

      if (this.localStream) {
        this.localStream.getTracks().forEach((track) => track.stop());
        this.localStream = null;
      }

      this.remoteStream = null;
      if (this.options.mediaEl) {
        this.options.mediaEl.srcObject = null;
      }

      if (updateState) {
        this.setState('stopped', '已停止', 'info');
      }
    }

    setState(state, message, tone) {
      this.state = state;
      window.GoLiveApp.setStatus(this.options.statusEl, message, tone);
      this.updateButtons();
    }

    updateButtons() {
      if (!this.options.startBtn) {
        return;
      }
      const busy = this.state === 'preparing' || this.state === 'connecting';
      const connected = this.state === 'connected';
      this.options.startBtn.disabled = busy || connected;
      this.options.startBtn.textContent = this.state === 'error'
        ? RETRY_LABELS[this.kind]
        : connected
          ? ACTIVE_LABELS[this.kind]
          : START_LABELS[this.kind];

      if (this.options.stopBtn) {
        this.options.stopBtn.disabled = !(busy || connected || this.state === 'error');
      }
    }

    log(message) {
      window.GoLiveApp.appendLog(this.options.logEl, message);
    }
  }

  function attachPageHideCleanup(session) {
    window.addEventListener('pagehide', () => {
      session.manualStop = true;
      session.cleanup();
    });
  }

  window.GoLiveWebRTC = {
    createPublisherSession(options) {
      const session = new WebRTCSession('publish', options);
      attachPageHideCleanup(session);
      return session;
    },
    createPlayerSession(options) {
      const session = new WebRTCSession('play', options);
      attachPageHideCleanup(session);
      return session;
    },
  };
})();
