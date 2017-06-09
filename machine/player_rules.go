package machine

import "log"

type PlayerRule interface {
	//检验条件
	Condition( player *Player ) bool
	//进行处理
	Action( player *Player )
	//进行处理
	AddPrompt( player *Player, action Action )
}

func interface_player_rule_add_prompt( player *Player, action Action ){
	player.HasPrompt = true
	player.PromptId += 1
	player.ActionDict[ player.PromptId ] = action

}

//==================================玩家暗杠====================================
type PlayerConcealedKongRule struct {}
func (self *PlayerConcealedKongRule) Condition( player *Player ) bool {
	if player.Table.CardsRest.Len() <= 0 {
		return false
	} else {
		record := make(map[int]int)
		for _,v := range player.CardsInHand {
			if v == 0 {
				continue
			} else {
				record[v]++
			}
		}
		log.Println(record)
		flag := false
		var action = Action{}
		for k,v := range record {
			if v >= 4 {
				action.ActionId = PlayerAction["PLAYER_ACTION_KONG_CONCEALED"]
				action.ActionCard = k
				action.Weight = PlayerAction["PLAYER_ACTION_KONG_CONCEALED"]
				self.AddPrompt( player, action )
				flag = true
			}
		}
		return flag
	}
}
func (self *PlayerConcealedKongRule) Action( player *Player ) {
	player.Machine.Trigger( &PlayerConcealedKongRuleState{} )
}
func (self *PlayerConcealedKongRule) AddPrompt( player *Player, action Action ) {
	interface_player_rule_add_prompt( player, action )
}


//==================================玩家过路杠====================================
type PlayerPongKongRule struct {}
func (self *PlayerPongKongRule) Condition( player *Player ) bool {
	if player.Table.CardsRest.Len() <= 0 || len(player.CardsPong) <= 0 {
		return false
	} else {
		flag := false
		var action = Action{}
		for _,v := range player.CardsPong {
			if v == 0 {
				continue
			} else {
				for _,card := range player.CardsInHand {
					if card == v {
						action.ActionId = PlayerAction["PLAYER_ACTION_KONG_PONG"]
						action.ActionCard = card
						action.Weight = PlayerAction["PLAYER_ACTION_KONG_PONG"]
						self.AddPrompt( player, action )
						flag = true
					}
				}
			}
		}
		return flag
	}
}
func (self *PlayerPongKongRule) Action( player *Player ) {
	player.Machine.Trigger( &PlayerPongKongRuleState{} )
}
func (self *PlayerPongKongRule) AddPrompt( player *Player, action Action ) {
	interface_player_rule_add_prompt( player, action )
}


//==================================玩家明杠====================================
type PlayerExposedKongRule struct {}
func (self *PlayerExposedKongRule) Condition( player *Player ) bool {
	if player.Table.CardsRest.Len() <= 0 {
		return false
	} else {
		//检测手中对应牌数量
		card_num := 0
		active_card := player.Table.ActiveCard
		for _,card := range player.CardsInHand {
			if card == active_card {
				card_num++
			}
		}
		//大于三张才能明杠
		if card_num >= 3 {
			action := Action{
				PlayerAction["PLAYER_ACTION_KONG_EXPOSED"],
				active_card,
				[]int32{},
				PlayerAction["PLAYER_ACTION_KONG_EXPOSED"],
			}
			self.AddPrompt( player, action )
			return true
		} else {
			return false
		}
	}
}
func (self *PlayerExposedKongRule) Action( player *Player ) {
	player.Machine.Trigger( &PlayerExposedKongRuleState{} )
}
func (self *PlayerExposedKongRule) AddPrompt( player *Player, action Action ) {
	interface_player_rule_add_prompt( player, action )
}


//==================================玩家碰牌====================================
type PlayerPongRule struct {}
func (self *PlayerPongRule) Condition( player *Player ) bool {
	active_card := player.Table.ActiveCard
	num := 0
	for _,v := range player.CardsInHand {
		if v == active_card {
			num++
		}
	}
	if num >= 2 {
		action := Action{
			PlayerAction["PLAYER_ACTION_PONG"],
			active_card,
			[]int32{},
			PlayerAction["PLAYER_ACTION_PONG"],
		}
		self.AddPrompt( player, action )
		return true
	} else {
		return false
	}

}
func (self *PlayerPongRule) Action( player *Player ) {
	player.Machine.Trigger( &PlayerPongRuleState{} )
}
func (self *PlayerPongRule) AddPrompt( player *Player, action Action ) {
	interface_player_rule_add_prompt( player, action )
}