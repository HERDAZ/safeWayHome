APIKEY="zdolX5bnPDYS60bWnpZAbWcHEnyU6V5C" # iaLB
#APIKEY="ZQQeYb2FBt4RMf73LbafMsZJhpMrxP0T" # kGtH
USERNAME="newUsername"
PASSWORD="newPassword"
PHONENB="23896712"
EMAIL="newEmail@gmail.cum"
LATITUDE="10.234567"
LONGITUDE="9.234567"
HOMELATITUDE="10.234567"
HOMELONGITUDE="9.234567"
FRIENDID="kGtH"

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
		\rTest for POST /amHome -d 'apikey':$APIKEY\n"
	curl 87.106.79.94:8447/amHome -d "{\"apikey\":\"$APIKEY\"}"
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

	all)
	./testApp.sh getPosition && ./testApp.sh getHome && ./testApp.sh getLogin && ./testApp.sh postHome && ./testApp.sh postPosition && ./testApp.sh postSignup
	;;

	*)
	echo "BAD ARGUMENT" $1
	echo "Possible args :"
	echo "postSignup"
	echo "postHome"
	echo "postPosition"
	echo "postAmHome"
	echo "getLogin"
	echo "getHome"
	echo "getPosition"
	echo "getIsHome"
	exit 1
	;;
esac
