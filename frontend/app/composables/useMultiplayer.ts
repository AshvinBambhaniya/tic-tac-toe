export const useMultiplayer = () => {
    const { $ws } = useNuxtApp();
    const { updateFromServer, isOpponentDisconnected } = useGameState();
    const config = useRuntimeConfig();
    const router = useRouter();
    
    // Track if we have an active listener to avoid duplicates if composable is reused
    let unregister: (() => void) | null = null;

    const connect = async (id: string, onMatchFound?: (gameId: string) => void) => {
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const host = config.public.apiUrl.replace(/^http(s)?:\/\//, '');
        const url = `${protocol}//${host}/api/v1/games/ws/${id}`;

        // Register message handler
        if (unregister) unregister();
        
        unregister = $ws.onMessage((msg: any) => {
            console.log('Message from server via service:', msg);
            switch (msg.type) {
                case 'STATE_UPDATE':
                    updateFromServer(msg.payload);
                    break;
                case 'MATCH_FOUND':
                    if (onMatchFound) {
                        onMatchFound(msg.payload.game_id);
                    } else {
                        router.push(`/game/${msg.payload.game_id}`);
                    }
                    break;
                case 'OPPONENT_LEFT':
                    isOpponentDisconnected.value = true;
                    break;
                case 'ERROR':
                    alert(msg.payload);
                    break;
            }
        });

        await $ws.connect(url);
    };

    const sendMove = (subGridIndex: number, cellIndex: number) => {
        $ws.send('MOVE', {
            subGridIndex,
            cellIndex
        });
    };

    const sendForfeit = () => {
        $ws.send('FORFEIT', {});
    };

    const disconnect = () => {
        if (unregister) {
            unregister();
            unregister = null;
        }
        $ws.disconnect();
    };

    return {
        connect,
        sendMove,
        sendForfeit,
        disconnect,
        readyState: computed(() => $ws.readyState)
    };
};
