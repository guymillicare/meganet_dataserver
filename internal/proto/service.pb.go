// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.0
// 	protoc        v5.26.1
// source: internal/proto/service.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// The request message containing the information needed to list games.
type ListPrematchRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ListPrematchRequest) Reset() {
	*x = ListPrematchRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_proto_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListPrematchRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListPrematchRequest) ProtoMessage() {}

func (x *ListPrematchRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListPrematchRequest.ProtoReflect.Descriptor instead.
func (*ListPrematchRequest) Descriptor() ([]byte, []int) {
	return file_internal_proto_service_proto_rawDescGZIP(), []int{0}
}

// The response message containing the list of games.
type ListPrematchResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data []*Prematch `protobuf:"bytes,1,rep,name=data,proto3" json:"data,omitempty"`
}

func (x *ListPrematchResponse) Reset() {
	*x = ListPrematchResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_proto_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListPrematchResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListPrematchResponse) ProtoMessage() {}

func (x *ListPrematchResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListPrematchResponse.ProtoReflect.Descriptor instead.
func (*ListPrematchResponse) Descriptor() ([]byte, []int) {
	return file_internal_proto_service_proto_rawDescGZIP(), []int{1}
}

func (x *ListPrematchResponse) GetData() []*Prematch {
	if x != nil {
		return x.Data
	}
	return nil
}

// A Prematch message represents a game in the sportsbook.
type Prematch struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id         string  `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	StartDate  string  `protobuf:"bytes,2,opt,name=start_date,json=startDate,proto3" json:"start_date,omitempty"`
	HomeTeam   string  `protobuf:"bytes,3,opt,name=home_team,json=homeTeam,proto3" json:"home_team,omitempty"`
	AwayTeam   string  `protobuf:"bytes,4,opt,name=away_team,json=awayTeam,proto3" json:"away_team,omitempty"`
	IsLive     bool    `protobuf:"varint,5,opt,name=is_live,json=isLive,proto3" json:"is_live,omitempty"`
	IsPopular  bool    `protobuf:"varint,6,opt,name=is_popular,json=isPopular,proto3" json:"is_popular,omitempty"`
	Tournament string  `protobuf:"bytes,7,opt,name=tournament,proto3" json:"tournament,omitempty"`
	Status     string  `protobuf:"bytes,8,opt,name=status,proto3" json:"status,omitempty"`
	Sport      string  `protobuf:"bytes,9,opt,name=sport,proto3" json:"sport,omitempty"`
	League     string  `protobuf:"bytes,10,opt,name=league,proto3" json:"league,omitempty"`
	Odds       []*Odds `protobuf:"bytes,11,rep,name=odds,proto3" json:"odds,omitempty"`
}

func (x *Prematch) Reset() {
	*x = Prematch{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_proto_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Prematch) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Prematch) ProtoMessage() {}

func (x *Prematch) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Prematch.ProtoReflect.Descriptor instead.
func (*Prematch) Descriptor() ([]byte, []int) {
	return file_internal_proto_service_proto_rawDescGZIP(), []int{2}
}

func (x *Prematch) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Prematch) GetStartDate() string {
	if x != nil {
		return x.StartDate
	}
	return ""
}

func (x *Prematch) GetHomeTeam() string {
	if x != nil {
		return x.HomeTeam
	}
	return ""
}

func (x *Prematch) GetAwayTeam() string {
	if x != nil {
		return x.AwayTeam
	}
	return ""
}

func (x *Prematch) GetIsLive() bool {
	if x != nil {
		return x.IsLive
	}
	return false
}

func (x *Prematch) GetIsPopular() bool {
	if x != nil {
		return x.IsPopular
	}
	return false
}

func (x *Prematch) GetTournament() string {
	if x != nil {
		return x.Tournament
	}
	return ""
}

func (x *Prematch) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *Prematch) GetSport() string {
	if x != nil {
		return x.Sport
	}
	return ""
}

func (x *Prematch) GetLeague() string {
	if x != nil {
		return x.League
	}
	return ""
}

func (x *Prematch) GetOdds() []*Odds {
	if x != nil {
		return x.Odds
	}
	return nil
}

type Odds struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id                  string  `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	SportsBookName      string  `protobuf:"bytes,2,opt,name=sports_book_name,json=sportsBookName,proto3" json:"sports_book_name,omitempty"`
	Name                string  `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	Price               float64 `protobuf:"fixed64,4,opt,name=price,proto3" json:"price,omitempty"`
	Timestamp           float64 `protobuf:"fixed64,5,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	BetPoints           float64 `protobuf:"fixed64,6,opt,name=bet_points,json=betPoints,proto3" json:"bet_points,omitempty"`
	IsMain              bool    `protobuf:"varint,7,opt,name=is_main,json=isMain,proto3" json:"is_main,omitempty"`
	IsLive              bool    `protobuf:"varint,8,opt,name=is_live,json=isLive,proto3" json:"is_live,omitempty"`
	MarketName          string  `protobuf:"bytes,9,opt,name=market_name,json=marketName,proto3" json:"market_name,omitempty"`
	Market              string  `protobuf:"bytes,10,opt,name=market,proto3" json:"market,omitempty"`
	HomeRotationNumber  float64 `protobuf:"fixed64,11,opt,name=home_rotation_number,json=homeRotationNumber,proto3" json:"home_rotation_number,omitempty"`
	AwayRotationNumber  float64 `protobuf:"fixed64,12,opt,name=away_rotation_number,json=awayRotationNumber,proto3" json:"away_rotation_number,omitempty"`
	DeepLinkUrl         string  `protobuf:"bytes,13,opt,name=deep_link_url,json=deepLinkUrl,proto3" json:"deep_link_url,omitempty"`
	PlayerId            string  `protobuf:"bytes,14,opt,name=player_id,json=playerId,proto3" json:"player_id,omitempty"`
	Selection           string  `protobuf:"bytes,15,opt,name=selection,proto3" json:"selection,omitempty"`
	NormalizedSelection string  `protobuf:"bytes,16,opt,name=normalized_selection,json=normalizedSelection,proto3" json:"normalized_selection,omitempty"`
	SelectionLine       string  `protobuf:"bytes,17,opt,name=selection_line,json=selectionLine,proto3" json:"selection_line,omitempty"`
	SelectionPoints     float64 `protobuf:"fixed64,18,opt,name=selection_points,json=selectionPoints,proto3" json:"selection_points,omitempty"`
}

func (x *Odds) Reset() {
	*x = Odds{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_proto_service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Odds) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Odds) ProtoMessage() {}

func (x *Odds) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_service_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Odds.ProtoReflect.Descriptor instead.
func (*Odds) Descriptor() ([]byte, []int) {
	return file_internal_proto_service_proto_rawDescGZIP(), []int{3}
}

func (x *Odds) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Odds) GetSportsBookName() string {
	if x != nil {
		return x.SportsBookName
	}
	return ""
}

func (x *Odds) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Odds) GetPrice() float64 {
	if x != nil {
		return x.Price
	}
	return 0
}

func (x *Odds) GetTimestamp() float64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *Odds) GetBetPoints() float64 {
	if x != nil {
		return x.BetPoints
	}
	return 0
}

func (x *Odds) GetIsMain() bool {
	if x != nil {
		return x.IsMain
	}
	return false
}

func (x *Odds) GetIsLive() bool {
	if x != nil {
		return x.IsLive
	}
	return false
}

func (x *Odds) GetMarketName() string {
	if x != nil {
		return x.MarketName
	}
	return ""
}

func (x *Odds) GetMarket() string {
	if x != nil {
		return x.Market
	}
	return ""
}

func (x *Odds) GetHomeRotationNumber() float64 {
	if x != nil {
		return x.HomeRotationNumber
	}
	return 0
}

func (x *Odds) GetAwayRotationNumber() float64 {
	if x != nil {
		return x.AwayRotationNumber
	}
	return 0
}

func (x *Odds) GetDeepLinkUrl() string {
	if x != nil {
		return x.DeepLinkUrl
	}
	return ""
}

func (x *Odds) GetPlayerId() string {
	if x != nil {
		return x.PlayerId
	}
	return ""
}

func (x *Odds) GetSelection() string {
	if x != nil {
		return x.Selection
	}
	return ""
}

func (x *Odds) GetNormalizedSelection() string {
	if x != nil {
		return x.NormalizedSelection
	}
	return ""
}

func (x *Odds) GetSelectionLine() string {
	if x != nil {
		return x.SelectionLine
	}
	return ""
}

func (x *Odds) GetSelectionPoints() float64 {
	if x != nil {
		return x.SelectionPoints
	}
	return 0
}

type LiveOddsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *LiveOddsRequest) Reset() {
	*x = LiveOddsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_proto_service_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LiveOddsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LiveOddsRequest) ProtoMessage() {}

func (x *LiveOddsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_service_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LiveOddsRequest.ProtoReflect.Descriptor instead.
func (*LiveOddsRequest) Descriptor() ([]byte, []int) {
	return file_internal_proto_service_proto_rawDescGZIP(), []int{4}
}

type LiveOddsData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data    []*Data `protobuf:"bytes,1,rep,name=data,proto3" json:"data,omitempty"`
	EntryId string  `protobuf:"bytes,2,opt,name=entry_id,json=entryId,proto3" json:"entry_id,omitempty"`
	Type    string  `protobuf:"bytes,3,opt,name=type,proto3" json:"type,omitempty"`
}

func (x *LiveOddsData) Reset() {
	*x = LiveOddsData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_proto_service_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LiveOddsData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LiveOddsData) ProtoMessage() {}

func (x *LiveOddsData) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_service_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LiveOddsData.ProtoReflect.Descriptor instead.
func (*LiveOddsData) Descriptor() ([]byte, []int) {
	return file_internal_proto_service_proto_rawDescGZIP(), []int{5}
}

func (x *LiveOddsData) GetData() []*Data {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *LiveOddsData) GetEntryId() string {
	if x != nil {
		return x.EntryId
	}
	return ""
}

func (x *LiveOddsData) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

type Data struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BetName         string  `protobuf:"bytes,1,opt,name=bet_name,json=betName,proto3" json:"bet_name,omitempty"`
	BetPoints       float64 `protobuf:"fixed64,2,opt,name=bet_points,json=betPoints,proto3" json:"bet_points,omitempty"`
	BetPrice        float64 `protobuf:"fixed64,3,opt,name=bet_price,json=betPrice,proto3" json:"bet_price,omitempty"`
	BetType         string  `protobuf:"bytes,4,opt,name=bet_type,json=betType,proto3" json:"bet_type,omitempty"`
	GameId          string  `protobuf:"bytes,5,opt,name=game_id,json=gameId,proto3" json:"game_id,omitempty"`
	Id              string  `protobuf:"bytes,6,opt,name=id,proto3" json:"id,omitempty"`
	IsLive          bool    `protobuf:"varint,7,opt,name=is_live,json=isLive,proto3" json:"is_live,omitempty"`
	IsMain          bool    `protobuf:"varint,8,opt,name=is_main,json=isMain,proto3" json:"is_main,omitempty"`
	League          string  `protobuf:"bytes,9,opt,name=league,proto3" json:"league,omitempty"`
	PlayerId        string  `protobuf:"bytes,10,opt,name=player_id,json=playerId,proto3" json:"player_id,omitempty"`
	Selection       string  `protobuf:"bytes,11,opt,name=selection,proto3" json:"selection,omitempty"`
	SelectionLine   string  `protobuf:"bytes,12,opt,name=selection_line,json=selectionLine,proto3" json:"selection_line,omitempty"`
	SelectionPoints string  `protobuf:"bytes,13,opt,name=selection_points,json=selectionPoints,proto3" json:"selection_points,omitempty"`
	Sport           string  `protobuf:"bytes,14,opt,name=sport,proto3" json:"sport,omitempty"`
	Sportsbook      string  `protobuf:"bytes,15,opt,name=sportsbook,proto3" json:"sportsbook,omitempty"`
	Timestamp       float64 `protobuf:"fixed64,16,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
}

func (x *Data) Reset() {
	*x = Data{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_proto_service_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Data) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Data) ProtoMessage() {}

func (x *Data) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_service_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Data.ProtoReflect.Descriptor instead.
func (*Data) Descriptor() ([]byte, []int) {
	return file_internal_proto_service_proto_rawDescGZIP(), []int{6}
}

func (x *Data) GetBetName() string {
	if x != nil {
		return x.BetName
	}
	return ""
}

func (x *Data) GetBetPoints() float64 {
	if x != nil {
		return x.BetPoints
	}
	return 0
}

func (x *Data) GetBetPrice() float64 {
	if x != nil {
		return x.BetPrice
	}
	return 0
}

func (x *Data) GetBetType() string {
	if x != nil {
		return x.BetType
	}
	return ""
}

func (x *Data) GetGameId() string {
	if x != nil {
		return x.GameId
	}
	return ""
}

func (x *Data) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Data) GetIsLive() bool {
	if x != nil {
		return x.IsLive
	}
	return false
}

func (x *Data) GetIsMain() bool {
	if x != nil {
		return x.IsMain
	}
	return false
}

func (x *Data) GetLeague() string {
	if x != nil {
		return x.League
	}
	return ""
}

func (x *Data) GetPlayerId() string {
	if x != nil {
		return x.PlayerId
	}
	return ""
}

func (x *Data) GetSelection() string {
	if x != nil {
		return x.Selection
	}
	return ""
}

func (x *Data) GetSelectionLine() string {
	if x != nil {
		return x.SelectionLine
	}
	return ""
}

func (x *Data) GetSelectionPoints() string {
	if x != nil {
		return x.SelectionPoints
	}
	return ""
}

func (x *Data) GetSport() string {
	if x != nil {
		return x.Sport
	}
	return ""
}

func (x *Data) GetSportsbook() string {
	if x != nil {
		return x.Sportsbook
	}
	return ""
}

func (x *Data) GetTimestamp() float64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

var File_internal_proto_service_proto protoreflect.FileDescriptor

var file_internal_proto_service_proto_rawDesc = []byte{
	0x0a, 0x1c, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x15, 0x0a, 0x13, 0x4c, 0x69, 0x73, 0x74, 0x50, 0x72, 0x65,
	0x6d, 0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x3b, 0x0a, 0x14,
	0x4c, 0x69, 0x73, 0x74, 0x50, 0x72, 0x65, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x23, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x72, 0x65, 0x6d, 0x61,
	0x74, 0x63, 0x68, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0xb2, 0x02, 0x0a, 0x08, 0x50, 0x72,
	0x65, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x74, 0x61, 0x72, 0x74, 0x5f,
	0x64, 0x61, 0x74, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x74, 0x61, 0x72,
	0x74, 0x44, 0x61, 0x74, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x68, 0x6f, 0x6d, 0x65, 0x5f, 0x74, 0x65,
	0x61, 0x6d, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x68, 0x6f, 0x6d, 0x65, 0x54, 0x65,
	0x61, 0x6d, 0x12, 0x1b, 0x0a, 0x09, 0x61, 0x77, 0x61, 0x79, 0x5f, 0x74, 0x65, 0x61, 0x6d, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x61, 0x77, 0x61, 0x79, 0x54, 0x65, 0x61, 0x6d, 0x12,
	0x17, 0x0a, 0x07, 0x69, 0x73, 0x5f, 0x6c, 0x69, 0x76, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x06, 0x69, 0x73, 0x4c, 0x69, 0x76, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x69, 0x73, 0x5f, 0x70,
	0x6f, 0x70, 0x75, 0x6c, 0x61, 0x72, 0x18, 0x06, 0x20, 0x01, 0x28, 0x08, 0x52, 0x09, 0x69, 0x73,
	0x50, 0x6f, 0x70, 0x75, 0x6c, 0x61, 0x72, 0x12, 0x1e, 0x0a, 0x0a, 0x74, 0x6f, 0x75, 0x72, 0x6e,
	0x61, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x74, 0x6f, 0x75,
	0x72, 0x6e, 0x61, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12,
	0x14, 0x0a, 0x05, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x73, 0x70, 0x6f, 0x72, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x6c, 0x65, 0x61, 0x67, 0x75, 0x65, 0x18,
	0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6c, 0x65, 0x61, 0x67, 0x75, 0x65, 0x12, 0x1f, 0x0a,
	0x04, 0x6f, 0x64, 0x64, 0x73, 0x18, 0x0b, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x4f, 0x64, 0x64, 0x73, 0x52, 0x04, 0x6f, 0x64, 0x64, 0x73, 0x22, 0xda,
	0x04, 0x0a, 0x04, 0x4f, 0x64, 0x64, 0x73, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x28, 0x0a, 0x10, 0x73, 0x70, 0x6f, 0x72, 0x74,
	0x73, 0x5f, 0x62, 0x6f, 0x6f, 0x6b, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0e, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x73, 0x42, 0x6f, 0x6f, 0x6b, 0x4e, 0x61, 0x6d,
	0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x70, 0x72, 0x69, 0x63, 0x65, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x01, 0x52, 0x05, 0x70, 0x72, 0x69, 0x63, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x74,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x05, 0x20, 0x01, 0x28, 0x01, 0x52, 0x09,
	0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x1d, 0x0a, 0x0a, 0x62, 0x65, 0x74,
	0x5f, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x18, 0x06, 0x20, 0x01, 0x28, 0x01, 0x52, 0x09, 0x62,
	0x65, 0x74, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x12, 0x17, 0x0a, 0x07, 0x69, 0x73, 0x5f, 0x6d,
	0x61, 0x69, 0x6e, 0x18, 0x07, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x69, 0x73, 0x4d, 0x61, 0x69,
	0x6e, 0x12, 0x17, 0x0a, 0x07, 0x69, 0x73, 0x5f, 0x6c, 0x69, 0x76, 0x65, 0x18, 0x08, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x06, 0x69, 0x73, 0x4c, 0x69, 0x76, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x6d, 0x61,
	0x72, 0x6b, 0x65, 0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0a, 0x6d, 0x61, 0x72, 0x6b, 0x65, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x6d,
	0x61, 0x72, 0x6b, 0x65, 0x74, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6d, 0x61, 0x72,
	0x6b, 0x65, 0x74, 0x12, 0x30, 0x0a, 0x14, 0x68, 0x6f, 0x6d, 0x65, 0x5f, 0x72, 0x6f, 0x74, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x0b, 0x20, 0x01, 0x28,
	0x01, 0x52, 0x12, 0x68, 0x6f, 0x6d, 0x65, 0x52, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x4e,
	0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x30, 0x0a, 0x14, 0x61, 0x77, 0x61, 0x79, 0x5f, 0x72, 0x6f,
	0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x0c, 0x20,
	0x01, 0x28, 0x01, 0x52, 0x12, 0x61, 0x77, 0x61, 0x79, 0x52, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x22, 0x0a, 0x0d, 0x64, 0x65, 0x65, 0x70, 0x5f,
	0x6c, 0x69, 0x6e, 0x6b, 0x5f, 0x75, 0x72, 0x6c, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b,
	0x64, 0x65, 0x65, 0x70, 0x4c, 0x69, 0x6e, 0x6b, 0x55, 0x72, 0x6c, 0x12, 0x1b, 0x0a, 0x09, 0x70,
	0x6c, 0x61, 0x79, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x49, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x73, 0x65, 0x6c, 0x65,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x0f, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x65, 0x6c,
	0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x31, 0x0a, 0x14, 0x6e, 0x6f, 0x72, 0x6d, 0x61, 0x6c,
	0x69, 0x7a, 0x65, 0x64, 0x5f, 0x73, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x10,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x13, 0x6e, 0x6f, 0x72, 0x6d, 0x61, 0x6c, 0x69, 0x7a, 0x65, 0x64,
	0x53, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x25, 0x0a, 0x0e, 0x73, 0x65, 0x6c,
	0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x6c, 0x69, 0x6e, 0x65, 0x18, 0x11, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0d, 0x73, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x4c, 0x69, 0x6e, 0x65,
	0x12, 0x29, 0x0a, 0x10, 0x73, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x70, 0x6f,
	0x69, 0x6e, 0x74, 0x73, 0x18, 0x12, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0f, 0x73, 0x65, 0x6c, 0x65,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x22, 0x11, 0x0a, 0x0f, 0x4c,
	0x69, 0x76, 0x65, 0x4f, 0x64, 0x64, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x5e,
	0x0a, 0x0c, 0x4c, 0x69, 0x76, 0x65, 0x4f, 0x64, 0x64, 0x73, 0x44, 0x61, 0x74, 0x61, 0x12, 0x1f,
	0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x12,
	0x19, 0x0a, 0x08, 0x65, 0x6e, 0x74, 0x72, 0x79, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x65, 0x6e, 0x74, 0x72, 0x79, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79,
	0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x22, 0xcc,
	0x03, 0x0a, 0x04, 0x44, 0x61, 0x74, 0x61, 0x12, 0x19, 0x0a, 0x08, 0x62, 0x65, 0x74, 0x5f, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x62, 0x65, 0x74, 0x4e, 0x61,
	0x6d, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x62, 0x65, 0x74, 0x5f, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x73,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x09, 0x62, 0x65, 0x74, 0x50, 0x6f, 0x69, 0x6e, 0x74,
	0x73, 0x12, 0x1b, 0x0a, 0x09, 0x62, 0x65, 0x74, 0x5f, 0x70, 0x72, 0x69, 0x63, 0x65, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x01, 0x52, 0x08, 0x62, 0x65, 0x74, 0x50, 0x72, 0x69, 0x63, 0x65, 0x12, 0x19,
	0x0a, 0x08, 0x62, 0x65, 0x74, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x62, 0x65, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x17, 0x0a, 0x07, 0x67, 0x61, 0x6d,
	0x65, 0x5f, 0x69, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x67, 0x61, 0x6d, 0x65,
	0x49, 0x64, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02,
	0x69, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x69, 0x73, 0x5f, 0x6c, 0x69, 0x76, 0x65, 0x18, 0x07, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x06, 0x69, 0x73, 0x4c, 0x69, 0x76, 0x65, 0x12, 0x17, 0x0a, 0x07, 0x69,
	0x73, 0x5f, 0x6d, 0x61, 0x69, 0x6e, 0x18, 0x08, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x69, 0x73,
	0x4d, 0x61, 0x69, 0x6e, 0x12, 0x16, 0x0a, 0x06, 0x6c, 0x65, 0x61, 0x67, 0x75, 0x65, 0x18, 0x09,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6c, 0x65, 0x61, 0x67, 0x75, 0x65, 0x12, 0x1b, 0x0a, 0x09,
	0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x49, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x73, 0x65, 0x6c,
	0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x65,
	0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x25, 0x0a, 0x0e, 0x73, 0x65, 0x6c, 0x65, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x6c, 0x69, 0x6e, 0x65, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0d, 0x73, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x4c, 0x69, 0x6e, 0x65, 0x12, 0x29,
	0x0a, 0x10, 0x73, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x70, 0x6f, 0x69, 0x6e,
	0x74, 0x73, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x73, 0x65, 0x6c, 0x65, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x70, 0x6f,
	0x72, 0x74, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x12,
	0x1e, 0x0a, 0x0a, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x73, 0x62, 0x6f, 0x6f, 0x6b, 0x18, 0x0f, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0a, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x73, 0x62, 0x6f, 0x6f, 0x6b, 0x12,
	0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x10, 0x20, 0x01,
	0x28, 0x01, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x32, 0x9b, 0x01,
	0x0a, 0x11, 0x53, 0x70, 0x6f, 0x72, 0x74, 0x73, 0x62, 0x6f, 0x6f, 0x6b, 0x53, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x12, 0x47, 0x0a, 0x0c, 0x4c, 0x69, 0x73, 0x74, 0x50, 0x72, 0x65, 0x6d, 0x61,
	0x74, 0x63, 0x68, 0x12, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4c, 0x69, 0x73, 0x74,
	0x50, 0x72, 0x65, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x1b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x50, 0x72, 0x65, 0x6d,
	0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3d, 0x0a, 0x0c,
	0x53, 0x65, 0x6e, 0x64, 0x4c, 0x69, 0x76, 0x65, 0x4f, 0x64, 0x64, 0x73, 0x12, 0x16, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4c, 0x69, 0x76, 0x65, 0x4f, 0x64, 0x64, 0x73, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x13, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4c, 0x69, 0x76,
	0x65, 0x4f, 0x64, 0x64, 0x73, 0x44, 0x61, 0x74, 0x61, 0x30, 0x01, 0x42, 0x10, 0x5a, 0x0e, 0x69,
	0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_internal_proto_service_proto_rawDescOnce sync.Once
	file_internal_proto_service_proto_rawDescData = file_internal_proto_service_proto_rawDesc
)

func file_internal_proto_service_proto_rawDescGZIP() []byte {
	file_internal_proto_service_proto_rawDescOnce.Do(func() {
		file_internal_proto_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_internal_proto_service_proto_rawDescData)
	})
	return file_internal_proto_service_proto_rawDescData
}

var file_internal_proto_service_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_internal_proto_service_proto_goTypes = []interface{}{
	(*ListPrematchRequest)(nil),  // 0: proto.ListPrematchRequest
	(*ListPrematchResponse)(nil), // 1: proto.ListPrematchResponse
	(*Prematch)(nil),             // 2: proto.Prematch
	(*Odds)(nil),                 // 3: proto.Odds
	(*LiveOddsRequest)(nil),      // 4: proto.LiveOddsRequest
	(*LiveOddsData)(nil),         // 5: proto.LiveOddsData
	(*Data)(nil),                 // 6: proto.Data
}
var file_internal_proto_service_proto_depIdxs = []int32{
	2, // 0: proto.ListPrematchResponse.data:type_name -> proto.Prematch
	3, // 1: proto.Prematch.odds:type_name -> proto.Odds
	6, // 2: proto.LiveOddsData.data:type_name -> proto.Data
	0, // 3: proto.SportsbookService.ListPrematch:input_type -> proto.ListPrematchRequest
	4, // 4: proto.SportsbookService.SendLiveOdds:input_type -> proto.LiveOddsRequest
	1, // 5: proto.SportsbookService.ListPrematch:output_type -> proto.ListPrematchResponse
	5, // 6: proto.SportsbookService.SendLiveOdds:output_type -> proto.LiveOddsData
	5, // [5:7] is the sub-list for method output_type
	3, // [3:5] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_internal_proto_service_proto_init() }
func file_internal_proto_service_proto_init() {
	if File_internal_proto_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_internal_proto_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListPrematchRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_proto_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListPrematchResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_proto_service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Prematch); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_proto_service_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Odds); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_proto_service_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LiveOddsRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_proto_service_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LiveOddsData); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_proto_service_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Data); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_internal_proto_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_internal_proto_service_proto_goTypes,
		DependencyIndexes: file_internal_proto_service_proto_depIdxs,
		MessageInfos:      file_internal_proto_service_proto_msgTypes,
	}.Build()
	File_internal_proto_service_proto = out.File
	file_internal_proto_service_proto_rawDesc = nil
	file_internal_proto_service_proto_goTypes = nil
	file_internal_proto_service_proto_depIdxs = nil
}
