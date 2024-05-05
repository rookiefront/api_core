package model

import "errors"

func GetModel(tableName string, isArray bool) (interface{}, error) {
	if isArray {
		tableName += "_arr"
	}
	hash := map[string]any{}
	hash["sys_user"] = &SysUser{}
	hash["sys_user_arr"] = []SysUser{}
	hash["sys_step"] = &SysStep{}
	hash["sys_step_arr"] = []SysStep{}
	hash["sys_role"] = &SysRole{}
	hash["sys_role_arr"] = []SysRole{}
	hash["sys_menu"] = &SysMenu{}
	hash["sys_menu_arr"] = []SysMenu{}
	hash["sys_dict"] = &SysDict{}
	hash["sys_dict_arr"] = []SysDict{}
	hash["sys_dict_item"] = &SysDictItem{}
	hash["sys_dict_item_arr"] = []SysDictItem{}
	hash["blind_date_member"] = &BlindDateMember{}
	hash["blind_date_member_arr"] = []BlindDateMember{}
	hash["blind_date_member_interest"] = &BlindDateMemberInterest{}
	hash["blind_date_member_interest_arr"] = []BlindDateMemberInterest{}
	hash["blind_date_member_message"] = &BlindDateMemberMessage{}
	hash["blind_date_member_message_arr"] = []BlindDateMemberMessage{}
	hash["blind_date_member_message_detail"] = &BlindDateMemberMessageDetail{}
	hash["blind_date_member_message_detail_arr"] = []BlindDateMemberMessageDetail{}
	hash["blind_date_ta"] = &BlindDateTa{}
	hash["blind_date_ta_arr"] = []BlindDateTa{}
	if _, ok := hash[tableName]; ok {
		return hash[tableName], nil
	}
	return "", errors.New("未找到")
}
