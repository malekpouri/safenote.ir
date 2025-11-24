// Type definitions for the WASM functions attached to window
declare global {
    interface Window {
        Go: any;
        encryptNote: (text: string, keyHex: string) => { encrypted?: string; error?: string };
        decryptNote: (encryptedBase64: string, keyHex: string) => { decrypted?: string; error?: string };
    }
}

export class EncryptionService {
    private wasmLoaded = false;

    async init() {
        if (this.wasmLoaded) return;

        if (!window.Go) {
            // Load the wasm_exec.js script dynamically if not present
            await this.loadScript('/wasm_exec.js');
        }

        const go = new window.Go();
        const result = await WebAssembly.instantiateStreaming(fetch('/safenote.wasm'), go.importObject);
        go.run(result.instance);
        this.wasmLoaded = true;
    }

    private loadScript(src: string): Promise<void> {
        return new Promise((resolve, reject) => {
            const script = document.createElement('script');
            script.src = src;
            script.onload = () => resolve();
            script.onerror = () => reject(new Error(`Failed to load script: ${src}`));
            document.body.appendChild(script);
        });
    }

    async generateKey(): Promise<string> {
        const key = new Uint8Array(32); // 256-bit key
        crypto.getRandomValues(key);
        return Array.from(key).map(b => b.toString(16).padStart(2, '0')).join('');
    }

    async encrypt(text: string, keyHex: string): Promise<string> {
        await this.init();
        const result = window.encryptNote(text, keyHex);
        if (result.error) {
            throw new Error(result.error);
        }
        return result.encrypted!;
    }

    async decrypt(encryptedBase64: string, keyHex: string): Promise<string> {
        await this.init();
        const result = window.decryptNote(encryptedBase64, keyHex);
        if (result.error) {
            throw new Error(result.error);
        }
        return result.decrypted!;
    }
}

export const encryptionService = new EncryptionService();
