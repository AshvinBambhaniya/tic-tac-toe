export type MessageHandler = (data: any) => void;

export class WebSocketService {
    private socket: WebSocket | null = null;
    private handlers: Set<MessageHandler> = new Set();
    private reconnectAttempts = 0;
    private maxReconnectAttempts = 5;
    private reconnectTimeout = 2000;

    connect(url: string): Promise<void> {
        return new Promise((resolve, reject) => {
            if (this.socket && this.socket.readyState === WebSocket.OPEN) {
                console.log('WebSocket already connected');
                resolve();
                return;
            }

            console.log(`Connecting to WebSocket: ${url}`);
            this.socket = new WebSocket(url);

            this.socket.onopen = () => {
                console.log('WebSocket connected successfully');
                this.reconnectAttempts = 0;
                resolve();
            };

            this.socket.onerror = (err) => {
                console.error('WebSocket Error:', err);
                reject(err);
            };

            this.socket.onmessage = async (event) => {
                let data = event.data;
                if (data instanceof Blob) {
                    data = await data.text();
                }

                try {
                    const parsedData = JSON.parse(data);
                    this.handlers.forEach(handler => handler(parsedData));
                } catch (err) {
                    console.error('Failed to parse WebSocket message:', err);
                }
            };

            this.socket.onclose = (event) => {
                console.log('WebSocket connection closed:', event.reason);
                this.socket = null;
                // Auto-reconnect logic could go here if desired
            };
        });
    }

    send(type: string, payload: any): void {
        if (this.socket && this.socket.readyState === WebSocket.OPEN) {
            this.socket.send(JSON.stringify({ type, payload }));
        } else {
            console.error('Cannot send message: WebSocket is not open', { type, payload });
        }
    }

    onMessage(handler: MessageHandler): () => void {
        this.handlers.add(handler);
        // Return an unregister function
        return () => this.handlers.delete(handler);
    }

    disconnect(): void {
        if (this.socket) {
            this.socket.close();
            this.socket = null;
        }
    }

    get readyState(): number {
        return this.socket?.readyState ?? WebSocket.CLOSED;
    }
}
