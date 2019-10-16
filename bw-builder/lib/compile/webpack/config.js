"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const WebpackBar = require("webpackbar");
const webpack = require("webpack");
const clean_webpack_plugin_1 = require("clean-webpack-plugin");
const util_1 = require("../util");
const context_1 = require("../update/context");
const vueLoader_1 = require("./vueLoader");
const uglifyjs_webpack_plugin_1 = require("uglifyjs-webpack-plugin");
class Config {
    constructor(option) {
        const context = new context_1.Context();
        const weexboxConfig = require(context.weexboxConfigPath);
        const env = option.env;
        const watch = option.watch;
        const watchOptions = {
            aggregateTimeout: 1000,
            ignored: /node_modules/,
        };
        const mode = 'none';
        const cache = watch ? true : false;
        const entry = util_1.Util.entries();
        const output = {
            path: context.distPath,
            filename: `${context.wwwDic}/[name].js`,
        };
        const resolve = {
            extensions: ['.mjs', '.js', '.vue', '.json'],
        };
        const plugins = [];
        if (watch === false) {
            plugins.push(new WebpackBar({
                name: 'WeexBox',
            }));
        }
        plugins.push(new clean_webpack_plugin_1.CleanWebpackPlugin(), new webpack.DefinePlugin({ 'process.env.NODE_ENV': JSON.stringify(env) }));
        if (env.toLowerCase().includes('release')) {
            plugins.push(new uglifyjs_webpack_plugin_1.UglifyJsPlugin({
                parallel: true,
            }));
        }
        if (watch === false) {
            plugins.push(new webpack.BannerPlugin({
                banner: '// { "framework": "Vue"} \n',
                raw: true,
                exclude: 'Vue',
            }));
        }
        const rules = [
            {
                test: /\.vue(\?[^?]+)?$/,
                use: [
                    {
                        loader: 'weex-loader',
                        options: vueLoader_1.vueLoader({ useVue: false }),
                    },
                ],
            },
            {
                test: /\.js$/,
                use: [
                    {
                        loader: 'babel-loader',
                    },
                ],
            },
            {
                test: /\.(png|jpg|gif)$/,
                use: [
                    {
                        loader: 'file-loader',
                        options: {
                            publicPath: weexboxConfig[env].imagePublicPath + '/static/',
                            name: '[name]_[hash].[ext]',
                            outputPath: context.staticDir,
                        },
                    },
                ],
            },
        ];
        const module = { rules };
        const node = context.nodeConfiguration;
        this.weexConfig = {
            watch,
            watchOptions,
            mode,
            cache,
            entry,
            output,
            resolve,
            plugins,
            module,
            node,
        };
    }
}
exports.Config = Config;
