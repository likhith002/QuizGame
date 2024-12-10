import {
  player,
  type Player,
} from "../net";

export class PlayerGame {
  joinPlayer(newPlayer: Player) {
    player.set(newPlayer);
  }
}