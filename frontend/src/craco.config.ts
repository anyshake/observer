import { readFileSync } from "fs";
import { resolve } from "path";
import { Configuration, DefinePlugin } from "webpack";

module.exports = {
	webpack: {
		configure: (webpackConfig: Configuration, { paths }: { paths: { appBuild: string } }) => {
			// Generate build tag with version and commit hash
			const getBuildTag = () => {
				let commit = "unknown";
				try {
					const filePath = resolve(__dirname, "..", "..", ".git", "logs", "HEAD");
					commit = readFileSync(filePath, "utf-8")
						.trim()
						.split("\n")
						.pop()!
						.split(" ")[1]
						.slice(0, 8);
				} catch {
					commit = "unknown";
				}
				let version = "custombuild";
				try {
					const filePath = resolve(__dirname, "..", "..", "VERSION");
					version = readFileSync(filePath, "utf-8").trim();
				} catch {
					version = "custombuild";
				}
				return JSON.stringify(`${version}-${commit}-${Math.floor(Date.now() / 1000)}`);
			};
			webpackConfig.plugins!.push(
				new DefinePlugin({
					"process.env.BUILD_TAG": getBuildTag()
				})
			);

			// Change build output path to ../dist
			const buildPath = resolve(__dirname, "..", "dist");
			webpackConfig.output!.path = buildPath;
			paths.appBuild = buildPath;

            // Add worker-loader for Web Workers
			webpackConfig.module!.rules!.unshift({
				test: /\.worker\.ts$/,
				use: {
					loader: "worker-loader",
					options: {
						// Use directory structure & typical names of chunks produces by "react-scripts"
						filename: "static/js/[name].[contenthash:8].js"
					}
				}
			});

			return webpackConfig;
		}
	}
};
