import { writable, type Writable } from "svelte/store";
import { GameState, NetService, PacketTypes, type ChangeGameState, type HostGamePacket, type Packet, type PlayerJoinPacket } from "../net";
import type { Player } from "../../model/quiz";


export const players:Writable<Player[]>=writable([])
export const gameState:Writable<GameState>=writable()
export const gameCode:Writable<string>=writable("")
export class HostGame{
    private net:NetService;

    constructor(){
        this.net=NetService.getInstance()

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
                    gameCode.set(data.gameCode)
                    break
                }

            case PacketTypes.ChangeState:
                {
                    let data=packet as ChangeGameState
                    gameState.set(data.state)
                    break
                }
        }

    }

}