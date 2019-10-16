"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const fs = require('fs');
const WebSocket = require("ws");
const EventEmitter = require("events");
const runner_1 = require("../common/runner");
const ws_1 = require("../server/ws");
class Runner extends EventEmitter {
    constructor(options) {
        super();
        this.config = options;
        this.on('error', e => {
            console.error(e);
        });
    }
    async startServer() {
        if (this.wsServer) {
            return this.wsServer;
        }
        const config = this.config;
        this.wsServer = new ws_1.default({
            //这个是deploy文件夹。
            staticFolder: config.jsBundleFolderPath,
        });
        console.log("config.jsBundleFolderPath---------------------",config.jsBundleFolderPath)
        await this.wsServer.init();
    }
    transmitEvent(outEvent) {
        outEvent.on(runner_1.messageType.outputError, message => {
            this.emit(runner_1.messageType.outputError, message);
        });
        outEvent.on(runner_1.messageType.outputLog, message => {
            this.emit(runner_1.messageType.outputLog, message);
        });
    }
    watchFileChange() {
        if (this.filesWatcher) {
            this.filesWatcher.close();
        }
        this.filesWatcher = fs.watch(this.config.jsBundleFolderPath, {
            recursive: true,
        }, (type, name) => {
            name = name.replace(/\\/g, '\/');
            const wsServer = this.wsServer;
            const serverInfo = wsServer.getServerInfo();
            const wsS = wsServer.getWs();
            if (!wsS) {
                return;
            }
            for (const ws of wsS) {
                if (ws.readyState === WebSocket.OPEN && name.includes('www')) {
                    ws.send(JSON.stringify({
                        method: 'WXReloadBundle',
                        params: `http://${serverInfo.hostname}:${serverInfo.port}/${name}`,
                    }));
                    console.log("给原生发送了什么东西====",`http://${serverInfo.hostname}:${serverInfo.port}/${name}`)
                }
            }
        });
        return true;
    }
    async run(options) {
        let appPath;
        try {
            this.emit(runner_1.messageType.state, runner_1.runnerState.start);
            await this.startServer();
            const serverInfo = this.wsServer.getServerInfo();
            this.emit(runner_1.messageType.state, runner_1.runnerState.startServerDone, `ws://${serverInfo.hostname}:${serverInfo.port}`);
            this.watchFileChange();
            this.emit(runner_1.messageType.state, runner_1.runnerState.watchFileChangeDone);
            this.emit(runner_1.messageType.state, runner_1.runnerState.done);
        }
        catch (error) {
            throw error;
        }
    }
    dispose() {
        this.filesWatcher.close();
        this.wsServer.dispose();
    }
}
exports.default = Runner;
