"use strict";

const express = require('express');
const { Wit } = require('node-wit');
const mongodb = require('mongodb');
const { co } = require('co');
const moment = require('moment')

const app = express();
const port = process.env.PORT || '80';
const host = process.env.HOST || '';
const dbAddr = process.env.DBADDR;
const witaiToken = process.env.WITAITOKEN;

const bodyParser = require('body-parser');
app.use(bodyParser.text({
	type: "text/plain"
}));


const MOST = 1;
const NOT = 2;
const NONE = 3;

if (!witaiToken) {
	console.error("please set WITAITOKEN to your wit.ai app token");
	process.exit(1);
}

if (!dbAddr){
	console.error("please set dbAddr to yourDatabase Address");
	process.exit(1);
}

const witaiClient = new Wit({ accessToken: witaiToken });

co(function*() {
	app.locals.db = yield mongodb.MongoClient.connect(`mongodb://${dbAddr}/info344`);
}).catch(function(err) {
	console.log(err.stack);
});

function whenHandler(req, res, data){
	co(function*(){
		var currentuser = JSON.parse(req.header('X-User'))
		var db = req.app.locals.db;
		var username = data.entities.user[0].value.toLowerCase();
		if(username == 'i' || username =='my'){
			username = currentuser.firstName.toLowerCase()
		}

		var users = (yield db.collection('users').find().toArray());
		var user;
		users.map((u)=>{
			if(u.firstname.toLowerCase() == username){
				user = u;
			}
		})

		if(!user){
			res.send(`I can't find the user.`)
			return;
		}

		//get list of all channels that the user is involved
		var channels = (yield db.collection('channels').find().toArray());

		channels = channels.filter((ch)=>{
			return (ch.members.indexOf(currentuser.id) != -1) && (ch.members.indexOf(user._id) != -1)
		});

		if(data.entities.channel){
			channels = channels.filter((ch)=>{
				return ch.name.toLowerCase() == data.entities.channel[0].value.toLowerCase()
			});	
		}

		if(!channels){
			res.send(`The user is not in your channels`)
			return;
		}
		var chanList = [];
		channels.map((ch)=>{
			chanList.push(ch._id)
		})

		var lastMsg = (yield db.collection('messages').find({creatorid: user._id}).sort({createdat:-1}).limit(1).toArray());
		var prettyDate = moment(lastMsg[0].createdat).format('LLL')
		var the = (data.entities.channel) ? data.entities.channel[0].value.toLowerCase() : 'the'


		if(lastMsg.length == 0){
			res.send(`No Message Found in ${the} channel`)
			return;
		} 
		res.send(`It was ${prettyDate}`)
		return;
	})
}

// "How many posts have I made to the XYZ channel?": 
// "How many posts did I make to the XYZ channel yesterday?" 
// "How many posts did I make to the XYZ channel on Monday?"
function howManyHandler(req, res, data){
	co(function*(){
		var currentuser = JSON.parse(req.header('X-User'))
		var db = req.app.locals.db;
		var username = data.entities.user[0].value.toLowerCase();
		if(username == 'i' || username =='my'){
			username = currentuser.firstName.toLowerCase()
		}

		var users = (yield db.collection('users').find().toArray());
		var user;

		users.map((u)=>{
			if(u.firstname.toLowerCase() == username){
				user = u;
			}
		})

		if(!user){
			res.send(`I can't find the user.`)
			return;
		}

		//get list of all channels that the user is involved
		var channels = (yield db.collection('channels').find().toArray());

		channels = channels.filter((ch)=>{
			return (ch.members.indexOf(currentuser.id) != -1) && (ch.members.indexOf(user._id) != -1)
		});

		if(data.entities.channel){
			channels = channels.filter((ch)=>{
				return ch.name.toLowerCase() == data.entities.channel[0].value.toLowerCase()
			});	
		}

		if(!channels){
			res.send(`The user is not in your channels`)
			return;
		}

		var chanList = [];
		channels.map((ch)=>{
			chanList.push(ch._id)
		})

		var messages = (yield db.collection('messages').find().toArray());
		if(!messages){
			res.send(`There is no message made`)
			return;
		}
		if(data.entities.datetime){
			messages = messages.filter((msg)=>{
				return moment(data.entities.datetime[0].value).format('MM/DD/YYYY') == moment(msg.createdat).format('MM/DD/YYYY')
			})
		}
		if(messages.length==0){
			res.send(`There is no message on ${moment(data.entities.datetime[0].value).format('LL')}`)
			return;
		}

		var count = 0;
		messages.map((msg)=>{
			if(chanList.indexOf(msg.channelid)!=-1){
				count++;
			}
		})
		res.send(`${username} made total of ${count} messages.`)
		return;
	});	
}


//"Who has made the most posts to the XYZ channel?"
//"Who is in the XYZ channel?"
//"Who hasn't posted to the XYZ channel?"
//"who has never posted to the XYZ channel?":
function whoHandler(req, res, data){
    if (data.entities.most) {
        handleMembers(req, res, data, MOST);
    } else if (data.entities.negative) {
        handleMembers(req, res, data, NOT);
    } else {
        handleMembers(req, res, data, NONE);
    }
}

function handleMembers(req, res, data, filter){
	co(function*(){
		var currentuser = JSON.parse(req.header('X-User'))
		var db = req.app.locals.db;

		var users = (yield db.collection('users').find().toArray());

		//get list of all channels that the user is involved
		var channels = (yield db.collection('channels').find().toArray());

		channels = channels.filter((ch)=>{
			return (ch.members.indexOf(currentuser.id) != -1)
		});

		if(data.entities.channel){
			channels = channels.filter((ch)=>{
				return ch.name.toLowerCase() == data.entities.channel[0].value.toLowerCase()
			});	
		}
		if(!channels){
			res.send(`cannot find the channel`)
			return;
		}

		var chanList = [];
		channels.map((ch)=>{
			chanList.push(ch._id)
		})

		var messages = (yield db.collection('messages').find().toArray());
		messages = messages.filter((msg)=>{
			return chanList.indexOf(msg.channelid) != -1
		})


		if(filter != NONE){
			var map = {};
			messages.map((msg)=>{
				if(!map[msg.creatorid]){
					map[msg.creatorid] = 0;
				}
				map[msg.creatorid]++
			})
			if(filter == MOST){
				var max = 0;
				var list = []
				Object.keys(map).forEach((key)=>{
					var value = map[key]
					if(max < value){
						max = value;
						list = [];
						list.push(key);
					} else if (max == value){
						list.push(key);
					}
				})
				var result = [];
				users.map((user)=>{
					if(list.indexOf(user._id) != -1){
						result.push(user.firstname)
					}
				})
				res.send(`${result.toString()} have/has send the most messages`)
			} else {
				var result = []
				users.map((user)=>{
					if(Object.keys(map).indexOf(user._id) == -1 ){
						result.push(user.firstname)
					}
				})
				if(result.length == 0){
					res.send(`Everyone has posted.`)
					return
				}
				res.send(`${result.toString()} have/has never send the any messages`)
			}	
		} else {
			list = [];
			channels.map((ch)=>{
				ch.members.map((member)=>{
					if(list.indexOf(member) == -1){
						list.push(member)
					}
				})
			})
			var result = []
			users.map((user)=>{
				if(list.indexOf(user._id) != -1){
					result.push(user.firstname)
				}
			})
			if(result.length == 1){
				res.send(`${result.toString()} is the only person in the channel`)
			} else {
				res.send(`${result.toString()} are in the channel`)
			}
		}

	})
}

app.post("/v1/bot", (req, res, next) => {
    let q = req.body;
    console.log(`user is asking ${q}`);
	witaiClient.message(q)
		.then(data => {
			switch (data.entities.intent[0].value.toLowerCase()) {
                case "when":
					whenHandler(req, res, data);
					break;
                case "how many":
					howManyHandler(req, res, data);
					break;
				case "who":
					whoHandler(req, res, data);
					break;
				case "list":
					handleMembers(req, res, data, NONE);
					break;
				default:
					res.send("Sorry, I'm not sure how to answer that. Please try again.");
			}
		})
		.catch(next);
});

app.listen(port, host, () => {
	console.log(`server is listening at http://${host}:${port}`);
});
