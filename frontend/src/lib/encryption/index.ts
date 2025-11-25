// Type definitions for the WASM functions attached to window
declare global {
    interface Window {
        Go: any;
        encryptNote: (text: string, shortKey: string, password?: string) => { encrypted?: string; error?: string };
        decryptNote: (encryptedBase64: string, shortKey: string, password?: string) => { decrypted?: string; error?: string };
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
        // Generate a 6-character random string (Base62-like)
        const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789";
        let result = "";
        const values = new Uint32Array(6);
        crypto.getRandomValues(values);
        for (let i = 0; i < 6; i++) {
            result += charset[values[i] % charset.length];
        }
        return result;
    }

    async encrypt(text: string, shortKey: string, password?: string): Promise<string> {
        await this.init();
        const result = window.encryptNote(text, shortKey, password || "");
        if (result.error) {
            throw new Error(result.error);
        }
        return result.encrypted!;
    }

    async decrypt(encryptedBase64: string, shortKey: string, password?: string): Promise<string> {
        await this.init();
        const result = window.decryptNote(encryptedBase64, shortKey, password || "");
        if (result.error) {
            throw new Error(result.error);
        }
        return result.decrypted!;
    }

    async hashPassword(password: string): Promise<string> {
        const msgBuffer = new TextEncoder().encode(password);
        const hashBuffer = await crypto.subtle.digest('SHA-256', msgBuffer);
        const hashArray = Array.from(new Uint8Array(hashBuffer));
        const hashHex = hashArray.map(b => b.toString(16).padStart(2, '0')).join('');
        return hashHex;
    }
}

export const encryptionService = new EncryptionService();
