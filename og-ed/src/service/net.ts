import { writable, type Writable } from "svelte/store";
import { PlayerGame } from "./player/player";

export enum PacketTypes {
  Connect,
  Host,
  ShowQuestion,
  ChangeState,
  PlayerJoin,
  StartGame,
  Tick,
  Coordinates,
  GameSettings,
  ChooseWord,
  SelectedWord,
  LevelResult,
  PlayerExit,
}

export enum GameState {
  Lobby,
  Play,
  UpdatePlayer,
  Wait,
  Reveal,
  End,
}

export interface Player {
  id: string;
  name: string;
  points: number;
  profile?: string;
}

export interface Packet {
  id: PacketTypes;
}

export interface HostGamePacket extends Packet {
  quizId: string;
}

export interface ChangeGameState extends Packet {
  state: GameState;
  payload: any;
}

export interface PlayerJoinPacket extends Packet {
  player: Player;
  gameCode: string;
}

export interface TickPacket extends Packet {
  tick: number;
}

export interface GameSettingsPacket extends Packet {
  players: Player[];
  coordinates: DrawPoint[];
  currentPlayer: Player;
}

export interface ChooseWordPacket extends Packet {
  words: string[];
}
export interface SelectedWord extends Packet {
  word: string;
}

export interface DrawPoint extends Packet {
  x1: number;
  y1: number;
  x2: number;
  y2: number;
  color: string;
  lineWidth: string;
}

export interface ResultPacket extends Packet {
  result: Player[];
  type: string;
}

export const state: Writable<GameState> = writable();
export const players: Writable<Player[]> = writable([]);
export const player: Writable<Player> = writable();
export const currentPlayer: Writable<Player> = writable();
export const gameState: Writable<GameState> = writable();
export class NetService {
  private webSocket!: WebSocket;
  private textDecoder: TextDecoder = new TextDecoder();
  private textEncoder: TextEncoder = new TextEncoder();
  private onPacketCallbacks: Array<(packet: any) => void> = [];
  private playerGame: PlayerGame = new PlayerGame();
  private static instance: NetService;
  connect() {
    this.webSocket = new WebSocket("ws://localhost:5001/ws");
    this.webSocket.onopen = () => {
      console.log("Connection opened");
      //   this.sendPacket({
      //     id: 0,
      //     code: "122",
      //     name: "test",
      //   });
    };

    this.webSocket.onclose = () => {
      console.log("Connection closed from server....");
    };

    this.webSocket.onmessage = async (event: MessageEvent) => {
      const arrayBuffer = await event.data.arrayBuffer();
      const bytes = new Uint8Array(arrayBuffer);
      const pId = bytes[0];

      const packet = JSON.parse(this.textDecoder.decode(bytes.subarray(1)));
      packet.id = pId;

      if (pId == PacketTypes.PlayerJoin) {
        const playerPacket = packet as PlayerJoinPacket;
        this.playerGame.joinPlayer(playerPacket.player);
      } else if (pId == PacketTypes.ChangeState) {
        const gameStatePacket = packet as ChangeGameState;
        state.set(gameStatePacket.state);

        if (gameStatePacket.state == GameState.UpdatePlayer) {
          const currPlayer = gameStatePacket.payload.player as Player;
          currentPlayer.set(currPlayer);
        }
      }

      console.log("received", packet);

      if (this.onPacketCallbacks.length > 0) {
        this.onPacketCallbacks.forEach((callback) => callback(packet));
      }
    };
  }

  close() {
    this.sendPacket({
      id: PacketTypes.PlayerExit,
    });
  }

  join(name: string) {
    let packet = {
      id: PacketTypes.Connect,
      name: name,
    };
    this.sendPacket(packet);
  }
  public static getInstance(): NetService {
    // If the instance does not exist, create one
    if (!NetService.instance) {
      NetService.instance = new NetService();
      this.instance.connect();
    }
    return NetService.instance;
  }

  onPacket(callback: (packet: Packet) => void) {
    this.onPacketCallbacks.push(callback);
  }

  sendPacket(packet: Packet) {
    const pId: number = packet.id;
    const pIdArr = new Uint8Array([pId]);
    const packetData = JSON.stringify(packet, (key, val) =>
      key == "id" ? undefined : val
    );
    const packetDataArr = this.textEncoder.encode(packetData);
    const finalArr = new Uint8Array(pIdArr.length + packetDataArr.length);
    finalArr.set(pIdArr);
    finalArr.set(packetDataArr, pIdArr.length);
    this.webSocket.send(finalArr);
  }
}
