import resolve from "rollup-plugin-node-resolve";
import copy from 'rollup-plugin-copy';

export default {
    input: "node_modules/@polymer/lit-element/lit-element.js",
    output: {
        file: "asset/lit-element.bundle.js",
        format: "es"
    },
    plugins: [
        resolve(),
        copy({
            "node_modules/@webcomponents/webcomponentsjs/webcomponents-bundle.js": "asset/webcomponents-bundle.js",
            verbose: true
        })
    ]
};