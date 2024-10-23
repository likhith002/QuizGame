
export enum PacketTypes{
  Connect,
  Host,
  ShowQuestion,
  ChnageState,
  PlayerJoin,
  StartGame,
  Tick
}


export enum GameState{

  Lobby,
  Play,
  Reveal,
  End
 
}


export interface ChangeGameState{
  state:GameState
}

interface Player{
  id:string
  name:string

}

export interface PlayerJoinPacket{
  player: Player

}

export interface TickPacket{
  tick:number
}

export class NetService {
  private webSocket!: WebSocket;
  private textDecoder: TextDecoder = new TextDecoder();
  private textEncoder: TextEncoder = new TextEncoder();
  private onPacketCallback?: (packet: any) => void;

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

    this.webSocket.onmessage = async (event: MessageEvent) => {

      const arrayBuffer = await event.data.arrayBuffer();
      const bytes = new Uint8Array(arrayBuffer);
      const pId = bytes[0];

      const packet = JSON.parse(this.textDecoder.decode(bytes.subarray(1)));
      packet.id=pId
      if (this.onPacketCallback) {
        this.onPacketCallback(packet);
      }

    }
  }

  onPacket(callback: (packet: any) => void) {
    this.onPacketCallback = callback;
  }

  sendPacket(packet: any) {
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
