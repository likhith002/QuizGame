import { writable, type Writable } from "svelte/store";
import { GameState, NetService, PacketTypes, type ChangeGameState, type Packet, type PlayerJoinPacket } from "../net";



export const state:Writable<GameState>=writable()
export class PlayerGame{

    private net:NetService;

    constructor(){
        this.net=NetService.getInstance()
        this.net.onPacket((p)=>this.onPacket(p))
    }

    join(name:string){

        let packet={
            id:PacketTypes.Connect,
            name:name
        }
        this.net.sendPacket(packet)

    }
    onPacket(p:Packet){
        switch(p.id){
        case PacketTypes.ChangeState:{
            let data=p as ChangeGameState
            state.set(data.state)
            return
        }
            
        }

    }
}