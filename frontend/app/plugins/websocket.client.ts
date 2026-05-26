import { WebSocketService } from "../services/WebSocketService";

export default defineNuxtPlugin(() => {
    const wsService = new WebSocketService();

    return {
        provide: {
            ws: wsService
        }
    };
});
