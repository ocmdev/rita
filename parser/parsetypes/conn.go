package parsetypes

import (
	"github.com/ocmdev/rita/config"
	"gopkg.in/mgo.v2/bson"
)

type (
	// Conn provides a data structure for bro's connection data
	Conn struct {
		// ID is the id coming out of mongodb
		ID bson.ObjectId `bson:"_id,omitempty"`
		// TimeStamp of this connection
		TimeStamp int64 `bson:"ts" bro:"ts" brotype:"time"`
		// Uid is the Unique Id for this connection (generated by Bro)
		UID string `bson:"uid" bro:"uid" brotype:"string"`
		// Source is the source address for this connection
		Source string `bson:"id_origin_h" bro:"id.orig_h" brotype:"addr"`
		// SourcePort is the source port of this connection
		SourcePort int `bson:"id_origin_p" bro:"id.orig_p" brotype:"port"`
		// Destination is the destination of the connection
		Destination string `bson:"id_resp_h" bro:"id.resp_h" brotype:"addr"`
		// DestinationPort is the port at the destination host
		DestinationPort int `bson:"id_resp_p" bro:"id.resp_p" brotype:"port"`
		// Proto is the string protocol identifier for this connection
		Proto string `bson:"proto" bro:"proto" brotype:"enum"`
		// Service describes the service of this connection if there was one
		Service string `bson:"service" bro:"service" brotype:"string"`
		// Duration is the floating point representation of connection length
		Duration float64 `bson:"duration" bro:"duration" brotype:"interval"`
		// OrigBytes is the byte count coming from the origin
		OrigBytes int64 `bson:"orig_bytes" bro:"orig_bytes" brotype:"count"`
		// RespBytes is the byte count coming in on response
		RespBytes int64 `bson:"resp_bytes" bro:"resp_bytes" brotype:"count"`
		// ConnState has data describing the state of a connection
		ConnState string `bson:"conn_state" bro:"conn_state" brotype:"string"`
		// LocalOrigin denotes that the connection originated locally
		LocalOrigin bool `bson:"local_orig" bro:"local_orig" brotype:"bool"`
		// LocalResponse denote that the connection responded locally
		LocalResponse bool `bson:"local_resp" bro:"local_resp" brotype:"bool"`
		// MissedBytes keeps a count of bytes missed
		MissedBytes int64 `bson:"missed_bytes" bro:"missed_bytes" brotype:"count"`
		// History is a string containing historical information
		History string `bson:"history"  bro:"history" brotype:"string"`
		// OrigPkts is a count of origin packets
		OrigPkts int64 `bson:"orig_pkts"  bro:"orig_pkts" brotype:"count"`
		// OrigIpBytes is another origin data count
		OrigIPBytes int64 `bson:"orig_ip_bytes" bro:"orig_ip_bytes" brotype:"count"`
		// RespPkts counts response packets
		RespPkts int64 `bson:"resp_pkts" bro:"resp_pkts" brotype:"count"`
		// RespIpBytes gives the bytecount of response data
		RespIPBytes int64 `bson:"resp_ip_bytes" bro:"resp_ip_bytes" brotype:"count"`
		// TunnelParents lists tunnel parents
		TunnelParents []string `bson:"tunnel_parents" bro:"tunnel_parents" brotype:"set[string]"`
	}
)

//TargetCollection returns the mongo collection this entry should be inserted
//into
func (in *Conn) TargetCollection(config *config.StructureTableCfg) string {
	return config.ConnTable
}

//Indices gives MongoDB indices that should be used with the collection
func (in *Conn) Indices() []string {
	return []string{"$hashed:id_origin_h", "$hashed:id_resp_h", "-duration", "ts"}
}

//Normalize pre processes this type of entry before it is imported by rita
func (in *Conn) Normalize() {}
