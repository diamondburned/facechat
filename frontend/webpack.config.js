const MiniCssExtractPlugin = require("mini-css-extract-plugin")
const HtmlWebpackPlugin = require("html-webpack-plugin")
const path = require("path")

const mode = process.env.NODE_ENV || "development"
const prod = mode === "production"

module.exports = {
	entry: "./src/main.js",
	resolve: {
		extensions: [".mjs", ".js", ".svelte"],
	},
	output: {
		path: path.resolve(__dirname, "./dist"),
		filename: "bundle.js",
		publicPath: "/",
	},
	module: {
		rules: [
			{
				test: /\.svelte$/,
				exclude: /node_modules/,
				use: {
					loader: "svelte-loader",
					options: {
						emitCss: true,
						hotReload: true,
					},
				},
			},
			{
				test: /\.css$/,
				use: [
					prod ? MiniCssExtractPlugin.loader : "style-loader",
					"css-loader",
				],
			},
		],
	},
	mode,
	plugins: [
		new MiniCssExtractPlugin({
			filename: "bundle.css",
		}),
		new HtmlWebpackPlugin({
			template: "./public/index.html",
			minify: prod
				? {
						collapseWhitespace: true,
				  }
				: false,
		}),
	],
	devtool: prod ? false : "source-map",
}
