#TODO avoid code repetition

APIKEY="Z2BDRpDhXBcEux6GqA7sRKuF7F9sRIWA" # iaLB
FRIENDAPIKEY="ZQQeYb2FBt4RMf73LbafMsZJhpMrxP0T" # kGtH
BADAPIKEY="AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"
USERNAME="newUsername"
PASSWORD="newPassword"
PHONENB="23896712"
EMAIL="newEmail@gmail.cum"
LATITUDE="10.234567"
LONGITUDE="9.234567"
HOMELATITUDE="10.234567"
HOMELONGITUDE="9.234567"
FRIENDID="kGtH"
PYTHONID="jupx"
PYTHONUSERNAME="python"
PYTHONEMAIL="python@python.py"
PYTHONPHONENB="932718"
PYTHONPASSWD="python"
PYTHONAPIKEY="USdaxOt7t1iSQrHriFYpQUhGpE5wzT4Z"

#USERNAME=$PYTHONUSERNAME
#PASSWORD=$PYTHONPASSWORD
#PHONENB=$PYTHONPHONENB
#EMAIL=$PYTHONEMAIL

case $1 in
	getPosition)
	printf "\r----------------------------------------\n \
		\rTest for GET /position -H 'apikey': $APIKEY 'friendID': $FRIENDID\n"
	curl 87.106.79.94:8447/position -G -H "apikey: $APIKEY" -H "friendID":"$FRIENDID"
	printf "\n----------------------------------------\n"
	;;

	getHome)
	printf "\r---------------------------------------\n \
		\rTest for GET /home -H 'apikey':$APIKEY 'friendID':$FRIENDID\n"
	curl 87.106.79.94:8447/home -G -H "apikey: $APIKEY" -H "friendID":"$FRIENDID"
	printf "\n----------------------------------------\n"
	;;

	getLogin)
	printf "\r---------------------------------------\n \
		\rTest for GET /login -H 'username':$USERNAME -H 'password':$PASSWORD\n"
	curl 87.106.79.94:8447/login -G -H "username: $USERNAME" -H "password: $PASSWORD"
	printf "\n----------------------------------------\n"
	;;

	postHome)
	printf "\r---------------------------------------\n \
		\rTest for POST /home -d 'apikey':$APIKEY, 'latitude':$HOMELATITUDE, 'longitude':$HOMELONGITUDE\n"
	curl 87.106.79.94:8447/home -d "{\"apikey\":\"$APIKEY\",\"latitude\":$HOMELATITUDE,\"longitude\":$HOMELONGITUDE}"
	printf "\n----------------------------------------\n"
	;;

	postPosition)
	printf "\r---------------------------------------\n \
		\rTest for POST /position -d 'apikey':$APIKEY, 'latitude':$LATITUDE, 'longitude':$LONGITUDE\n"
	curl 87.106.79.94:8447/position -d '{"apikey":"'$APIKEY'","latitude":'$LATITUDE',"longitude":'$LONGITUDE'}'
	printf "\n----------------------------------------\n"
	;;

	postSignup)
	printf "\r---------------------------------------\n \
		\rTest for POST /signup -d 'username':$USERNAME, 'phoneNb':$PHONENB, 'email':$EMAIL, 'password':$PASSWORD\n"
	curl 87.106.79.94:8447/signup -d "{\"username\":\"$USERNAME\",\"phoneNb\":\"$PHONENB\",\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\"}"
	echo 
	printf "\n----------------------------------------\n"
	;;

	postAmHome)
	printf "\r---------------------------------------\n \
		\rTest for POST /amHome -d 'apikey':$FRIENDAPIKEY\n"
	curl 87.106.79.94:8447/amHome -d "{\"apikey\":\"$FRIENDAPIKEY\"}"
	echo 
	printf "\n----------------------------------------\n"
	;;

	getIsHome)
	printf "\r---------------------------------------\n \
		\rTest for GET /isHome -H 'apikey':$APIKEY -H 'friendID':$FRIENDID\n"
	curl 87.106.79.94:8447/isHome -G -H "apikey: $APIKEY" -H "friendID: $FRIENDID"
	printf "\n----------------------------------------\n"
	;;

	postCleanIsHome)
	printf "\r---------------------------------------\n \
		\rTest for POST /cleanIsHome\n"
	curl 87.106.79.94:8447/cleanIsHome -d {}
	echo
	printf "\n----------------------------------------\n"
	;;

	postInEmergency)
	printf "\r---------------------------------------\n \
		\rTest for POST /InEmerency -d 'apikey':$FRIENDAPIKEY\n"
	curl 87.106.79.94:8447/inEmergency -d "{\"apikey\":\"$FRIENDAPIKEY\"}"
	echo 
	printf "\n----------------------------------------\n"
	;;

	getInEmergency)
	printf "\r---------------------------------------\n \
		\rTest for GET /InEmerency -H 'apikey':$APIKEY -H 'friendID':$FRIENDID\n"
	curl 87.106.79.94:8447/inEmergency -G -H "apikey: $APIKEY" -H "friendID: $FRIENDID"
	printf "\n----------------------------------------\n"
	;;

	postCleanInEmergency)
	printf "\r---------------------------------------\n \
		\rTest for POST /cleanInEmerency\n"
	curl 87.106.79.94:8447/cleanInEmergency -d {}
	echo
	printf "\n----------------------------------------\n"
	;;

	postAddFriend)
	printf "\r---------------------------------------\n \
		\rTest for POST /addFriend -d 'apikey':$APIKEY -d 'friendID':$FRIENDID\n"
	curl 87.106.79.94:8447/addFriend -d "{\"apikey\":\"$APIKEY\",\"friendID\":\"$FRIENDID\"}"

	echo 
	printf "\n----------------------------------------\n"
	;;

	getFriends)
	printf "\r---------------------------------------\n \
		\rTest for GET /getFriend -H 'apikey':$APIKEY\n"
	curl 87.106.79.94:8447/getFriend -G -H "apikey: $APIKEY"
	echo
	printf "\n---------------------------------------------\n"
	;;

	fuzzBadAPIkey)
	printf "\r----------------------------------------\n \
		\rTest for GET /position -H 'apikey': $BADAPIKEY 'friendID': $FRIENDID\n"
		curl 87.106.79.94:8447/position -G -H "apikey: $BADAPIKEY" -H "friendID":"$FRIENDID"
	printf "\n----------------------------------------\n"

	printf "\r---------------------------------------\n \
		\rTest for GET /home -H 'apikey':$BADAPIKEY 'friendID':$FRIENDID\n"
		curl 87.106.79.94:8447/home -G -H "apikey: $BADAPIKEY" -H "friendID":"$FRIENDID"
	printf "\n----------------------------------------\n"
	
	printf "\r---------------------------------------\n \
		\rTest for POST /home -d 'apikey':$BADAPIKEY, 'latitude':$HOMELATITUDE, 'longitude':$HOMELONGITUDE\n"
		curl 87.106.79.94:8447/home -d "{\"apikey\":\"$BADAPIKEY\",\"latitude\":$HOMELATITUDE,\"longitude\":$HOMELONGITUDE}"
	printf "\n----------------------------------------\n"

	printf "\r---------------------------------------\n \
		\rTest for POST /position -d 'apikey':$BADAPIKEY, 'latitude':$LATITUDE, 'longitude':$LONGITUDE\n"
		curl 87.106.79.94:8447/position -d '{"apikey":"'$BADAPIKEY'","latitude":'$LATITUDE',"longitude":'$LONGITUDE'}'
	printf "\n----------------------------------------\n"

	printf "\r---------------------------------------\n \
		\rTest for POST /amHome -d 'apikey':$BADAPIKEY\n"
		curl 87.106.79.94:8447/amHome -d "{\"apikey\":\"$BADAPIKEY\"}"
	printf "\n----------------------------------------\n"

	printf "\r---------------------------------------\n \
		\rTest for GET /isHome -H 'apikey':$BADAPIKEY -H 'friendID':$FRIENDID\n"
		curl 87.106.79.94:8447/isHome -G -H "apikey: $BADAPIKEY" -H "friendID: $FRIENDID"
	printf "\n----------------------------------------\n"

	printf "\r---------------------------------------\n \
		\rTest for POST /InEmerency -d 'apikey':$BADAPIKEY\n"
		curl 87.106.79.94:8447/inEmergency -d "{\"apikey\":\"$BADAPIKEY\"}"
	printf "\n----------------------------------------\n"

	printf "\r---------------------------------------\n \
		\rTest for GET /InEmerency -H 'apikey':$BADAPIKEY -H 'friendID':$FRIENDID\n"
		curl 87.106.79.94:8447/inEmergency -G -H "apikey: $BADAPIKEY" -H "friendID: $FRIENDID"
	printf "\n----------------------------------------\n"

	printf "\r---------------------------------------\n \
		\rTest for POST /addFriend -d 'apikey':$BADAPIKEY -d 'friendID':$FRIENDID\n"
		curl 87.106.79.94:8447/addFriend -d "{\"apikey\":\"$BADAPIKEY\",\"friendID\":\"$FRIENDID\"}"
	printf "\n----------------------------------------\n"

	printf "\r---------------------------------------\n \
		\rTest for GET /getFriend -H 'apikey':$BADAPIKEY\n"
		curl 87.106.79.94:8447/getFriend -G -H "apikey: $BADAPIKEY"
	printf "\n---------------------------------------------\n"
	;;


	all)
	echo "WARNING : THIS WILL RESET ALL API KEYS !"
	#todo add input support instead
	sleep 10
	./testApp.sh getPosition && ./testApp.sh getHome && ./testApp.sh getLogin && ./testApp.sh postHome && ./testApp.sh postPosition && ./testApp.sh postSignup  && ./testApp.sh postCleanIsHome && ./testApp.sh getIsHome && ./testApp.sh postAmHome && ./testApp.sh getIsHome && ./testApp.sh postCleanInEmergency && ./testApp.sh getInEmergency && ./testApp.sh postInEmergency && ./testApp.sh getInEmergency && ./testApp.sh postAddFriend
	;;

	*)
	echo "BAD ARGUMENT" $1
	echo "Possible args :"
	echo ""
	echo "getLogin"
	echo "postSignup"
	echo "---------------------"
	echo "getPosition"
	echo "postPosition"
	echo "---------------------"
	echo "postHome"
	echo "getHome"
	echo "---------------------"
	echo "postAmHome"
	echo "getIsHome"
	echo "postCleanIsHome"
	echo "---------------------"
	echo "postInEmergency"
	echo "getInEmergency"
	echo "postCleanInEmergency"
	echo "---------------------"
	echo "fuzzBadAPIkey"
	exit 1
	;;
esac
