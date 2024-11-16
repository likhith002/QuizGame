import { writable, type Writable } from "svelte/store";
import { GameState, NetService, PacketTypes, type ChangeGameState, type HostGamePacket, type Packet, type PlayerJoinPacket } from "../net";
import type { Player } from "../../model/quiz";


export const players:Writable<Player[]>=writable([])
export const gameState:Writable<GameState>=writable(GameState.Lobby)
export class HostGame{
    private net:NetService;

    constructor(){
        this.net=new NetService();
        this.net.connect()
        this.net.onPacket(p=>this.onPacket(p))

    }

    hostQuiz(quizId:string){

        let packet:HostGamePacket={
            id:PacketTypes.Host,
            quizId:quizId
        }

        this.net.sendPacket(packet)
    }

    start(){
        this.net.sendPacket({
            id:PacketTypes.StartGame
        })

    }

    onPacket(packet:Packet){
        switch(packet.id){


            case PacketTypes.PlayerJoin:
                {

                    let data=packet as PlayerJoinPacket
                    players.update((p)=>[...p,data.player])
                    break
                }

            case PacketTypes.ChnageState:
                {
                    let data=packet as ChangeGameState
                    gameState.set(data.state)
                    break
                }
        }

    }

}