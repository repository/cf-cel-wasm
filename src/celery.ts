import "./wasm_exec.js";
import wasmModule from "./main.wasm";

export type CeleryResult<T = undefined> = T extends undefined ? { error?: string } : T & { error?: string };

export type Celery = {
	eval: (expression: string, inputData: Record<string, unknown>) => CeleryResult<{ result?: unknown }>;
	analyzeType: (expression: string, variableTypes: Record<string, unknown>) => CeleryResult<{ resultType?: string; isValid?: boolean }>;
	analyzeTypeUnknown: (
		expression: string,
		variableNames: string[] | Record<string, unknown>,
	) => CeleryResult<{ resultType?: string; isValid?: boolean }>;
};

declare global {
	var $__celery: Celery;
}

const go = new Go();
const instance = WebAssembly.instantiate(wasmModule, go.importObject);

let instantiated = false;

export async function getCelery(): Promise<Celery> {
	if (!instantiated) {
		go.run(await instance);
		instantiated = true;
	}

	return globalThis.$__celery;
}
