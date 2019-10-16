"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const fse = require("fs-extra");
const path = require("path");
const glob = require("glob");
const const_1 = require("./const");
class PlatformConfigResolver {
    constructor(def) {
        this.def = null;
        this.replacer = {
            plist(source, key, value) {
                const r = new RegExp('(<key>' + key + '</key>\\s*<string>)[^<>]*?</string>', 'g');
                if (key === 'WXEntryBundleURL' || key === 'WXSocketConnectionURL') {
                    if (key === 'WXEntryBundleURL') {
                        value = path.join('bundlejs', value);
                    }
                    if (!r.test(source)) {
                        return source.replace(/<\/dict>\n?\W*?<\/plist>\W*?\n?\W*?\n?$/i, match => `  <key>${key}</key>\n  <string>${value}</string>\n${match}`);
                    }
                }
                return source.replace(r, '$1' + value + '</string>');
            },
            xmlTag(source, key, value, tagName = 'string') {
                const r = new RegExp(`<${tagName} name="${key}" .*>[^<]+?</${tagName}>`, 'g');
                return source.replace(r, `<${tagName} name="${key}">${value}</${tagName}>`);
            },
            xmlAttr(source, key, value, tagName = 'string') {
                const r = new RegExp(`<${tagName} name="${key}"\\s* value="[^"]*?"\\s*/>`, 'g');
                return source.replace(r, `<${tagName} name="${key}" value="${value}"/>`);
            },
            regexp(source, regexp, value) {
                return source.replace(regexp, function (m, a, b) {
                    return a + value + (b || '');
                });
            },
            moveAndReplacePackageName(oldname, newname, basePath) {
                const oldPath = oldname.split('.').join('/');
                const newPath = newname.split('.').join('/');
                const javaSourcePath = 'app/src/main';
                const options = {};
                const files = glob.sync(path.join(basePath, javaSourcePath, '**/*.+(java|xml)'), options);
                if (Array.isArray(files)) {
                    files.forEach(file => {
                        let data = fse.readFileSync(file, 'utf8');
                        data = data.replace(new RegExp(oldname, 'ig'), newname);
                        fse.outputFileSync(file.replace(new RegExp(oldPath, 'ig'), newPath), data);
                    });
                }
                if (oldPath !== newPath && Array.isArray(files) && files.length > 0) {
                    fse.removeSync(path.join(basePath, javaSourcePath, 'java', oldPath));
                }
            },
        };
        this.resolveConfigDef = (source, configDef, config, key, basePath) => {
            if (configDef.type) {
                if (config[key] === undefined) {
                    console.warn('Config:[' + key + '] must have a value!');
                    return source;
                }
                return this.replacer[configDef.type](source, configDef.key, config[key]);
            }
            else {
                return configDef.handler(source, config[key], this.replacer, basePath);
            }
        };
        this.def = def;
    }
    resolve(config, basePath) {
        basePath = basePath || path.join(process.cwd(), `platforms/${config.platform}`);
        for (let d in this.def) {
            if (this.def.hasOwnProperty(d)) {
                const targetPath = path.join(basePath, d);
                let source = fse.readFileSync(targetPath).toString();
                for (let key in this.def[d]) {
                    if (this.def[d].hasOwnProperty(key)) {
                        const configDef = this.def[d][key];
                        if (Array.isArray(configDef)) {
                            configDef.forEach(def => {
                                source = this.resolveConfigDef(source, def, config, key, basePath);
                            });
                        }
                        else {
                            source = this.resolveConfigDef(source, configDef, config, key, basePath);
                        }
                    }
                }
                fse.writeFileSync(targetPath, source);
            }
        }
    }
}
const androidConfigResolver = new PlatformConfigResolver({
    'app/build.gradle': {
        AppId: {
            type: 'regexp',
            key: /(applicationId ")[^"]*(")/g,
        },
    },
    'app/src/main/res/values/strings.xml': {
        AppName: {
            type: 'xmlTag',
            key: 'app_name',
        },
    },
    'app/src/main/AndroidManifest.xml': {
        AppId: {
            handler: function (source, value, replacer, basePath) {
                if (!value) {
                    return source;
                }
                if (/package="(.*)"/.test(source)) {
                    let match = /package="(.*)"/.exec(source);
                    if (match[1]) {
                        replacer.moveAndReplacePackageName(match[1], value, basePath);
                        return source.replace(new RegExp(`${match[1]}`, 'ig'), value);
                    }
                    return source;
                }
                else {
                    return source;
                }
            },
        },
    },
});
const iOSConfigResolver = new PlatformConfigResolver({
    'WeexBoxPlayground/Info.plist': {
        AppName: {
            type: 'plist',
            key: 'CFBundleDisplayName',
        },
        Version: {
            type: 'plist',
            key: 'CFBundleShortVersionString',
        },
        BuildVersion: {
            type: 'plist',
            key: 'CFBundleVersion',
        },
        AppId: {
            type: 'plist',
            key: 'CFBundleIdentifier',
        },
        WeexBundle: {
            type: 'plist',
            key: 'WXEntryBundleURL',
        },
        Ws: {
            type: 'plist',
            key: 'WXSocketConnectionURL',
        },
    },
    'WeexBoxPlayground.xcodeproj/project.pbxproj': {
        CodeSign: [
            {
                type: 'regexp',
                key: /("?CODE_SIGN_IDENTITY(?:\[sdk=iphoneos\*])?"?\s*=\s*").*?(")/g,
            },
            {
                type: 'plist',
                key: 'CODE_SIGN_IDENTITY(\\[sdk=iphoneos\\*])?',
            },
        ],
        Profile: [
            {
                type: 'regexp',
                key: /(PROVISIONING_PROFILE\s*=\s*")[^"]*?(")/g,
            },
            {
                type: 'plist',
                key: 'PROVISIONING_PROFILE',
            },
        ],
    },
});
exports.default = {
    [const_1.PLATFORM_TYPES.ios]: iOSConfigResolver,
    [const_1.PLATFORM_TYPES.android]: androidConfigResolver,
};
