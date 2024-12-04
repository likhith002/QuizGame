
export enum PacketTypes{
  Connect,
  Host,
  ShowQuestion,
  ChangeState,
  PlayerJoin,
  StartGame,
  Tick,
  Coordinates,
  GameSettings
}


export enum GameState{
  Lobby,
  Play,
  Reveal,
  End
 
}






interface Player{
  id:string
  name:string

}



export interface Packet{
  id:PacketTypes
}



export interface HostGamePacket extends Packet{
  quizId:string
}

export interface ChangeGameState extends Packet{
  state:GameState
}


export interface PlayerJoinPacket extends Packet{
  player: Player
  gameCode:string

}

export interface TickPacket extends Packet{
  tick:number
}

export interface GameSettingsPacket extends Packet{
  players:Player[]
  coordinates:DrawPoint[]
}

export interface DrawPoint extends Packet{
  x1 :number    
	y1 :number    
	x2 :number   
	y2 : number    
	color :string 
	lineWidth:string 

}

export class NetService {
  private webSocket!: WebSocket;
  private textDecoder: TextDecoder = new TextDecoder();
  private textEncoder: TextEncoder = new TextEncoder();
  private onPacketCallback?: (packet: any) => void;
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

    this.webSocket.onclose=()=>{
      console.log("Connection closed from server....")
    }

    this.webSocket.onmessage = async (event: MessageEvent) => {


      const arrayBuffer = await event.data.arrayBuffer();
      const bytes = new Uint8Array(arrayBuffer);
      const pId = bytes[0];

      const packet = JSON.parse(this.textDecoder.decode(bytes.subarray(1)));
      packet.id=pId
      if (this.onPacketCallback) {
        console.log("Received packet",packet)
        this.onPacketCallback(packet);
      }

    }
  }

  public static getInstance(): NetService {
    // If the instance does not exist, create one
    if (!NetService.instance) {
      NetService.instance = new NetService();
      this.instance.connect()
    }
    return NetService.instance;
  }

  onPacket(callback: (packet: Packet) => void) {
    this.onPacketCallback = callback;
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
