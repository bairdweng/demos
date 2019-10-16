"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
var runnerState;
(function (runnerState) {
    runnerState[runnerState["start"] = 0] = "start";
    runnerState[runnerState["startServerDone"] = 1] = "startServerDone";
    runnerState[runnerState["setNativeConfigDone"] = 2] = "setNativeConfigDone";
    runnerState[runnerState["copyJsBundleDone"] = 3] = "copyJsBundleDone";
    runnerState[runnerState["watchFileChangeDone"] = 4] = "watchFileChangeDone";
    runnerState[runnerState["buildNativeDone"] = 5] = "buildNativeDone";
    runnerState[runnerState["installAndLaunchAppDone"] = 6] = "installAndLaunchAppDone";
    runnerState[runnerState["done"] = 7] = "done";
})(runnerState = exports.runnerState || (exports.runnerState = {}));
var messageType;
(function (messageType) {
    messageType["state"] = "state";
    messageType["outputLog"] = "outputLog";
    messageType["outputError"] = "outputError";
})(messageType = exports.messageType || (exports.messageType = {}));
