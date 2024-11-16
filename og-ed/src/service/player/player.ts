import { NetService, PacketTypes, type Packet } from "../net";



export class PlayerGame{

    private net:NetService;
    constructor(){
        this.net=new NetService()
        this.net.connect()
        this.net.onPacket((p)=>this.onPacket(p))
    }

    join(name:string,code:string){

        let packet={
            id:PacketTypes.Connect,
            code:code,
            name:name
        }

        this.net.sendPacket(packet)

    }
    onPacket(p:Packet){


    }
}