APIKEY="TEST"
USERNAME="newUsername"
PASSWORD="newPassword"
PHONENB="23896712"
EMAIL="newEmail@gmail.cum"
LATITUDE="10.234567"
LONGITUDE="9.234567"
FRIENDID="PYTH"

case $1 in
	getPosition)
	printf "\r----------------------------------------\n \
		\rTest for GET /position -H 'apikey': $APIKEY 'userID': $FRIENDID\n"
	curl 87.106.79.94:8447/position -G -H "apikey: $APIKEY" -H "userID":"$FRIENDID"
	printf "\n----------------------------------------\n"
	;;

	getHome)
	printf "\r---------------------------------------\n \
		\rTest for GET /home -H 'apikey':$APIKEY 'userID':$FRIENDID\n"
	curl 87.106.79.94:8447/home -G -H "apikey: $APIKEY" -H "userID":"$FRIENDID"
	printf "\n----------------------------------------\n"
	;;

	getLogin)
	printf "\r---------------------------------------\n \
		\rTest for GET /login -H 'username':$USERNAME -H 'password':$PASSWORD\n"
	curl 87.106.79.94:8447/login -G -H "username: $USERNAME" -H "password: $PASSWORD"
	printf "\n----------------------------------------\n"
	;;

	postLogin)
	printf "\r---------------------------------------\n \
		\rTest for POST /login -d 'apikey':$APIKEY, 'latitude':$LATITUDE, 'longitude':$LONGITUDE\n"
	curl 87.106.79.94:8447/home -d "{'apikey':"$APIKEY",'latitude':"$LATITUDE",'longitude':"$LONGITUDE"}"
	printf "\n----------------------------------------\n"
	;;

	postPosition)
	printf "\r---------------------------------------\n \
		\rTest for POST /position -d 'apikey':$APIKEY, 'latitude':$LATITUDE, 'longitude':$LONGITUDE\n"
	curl 87.106.79.94:8447/position -d '{"apikey":"'$APIKEY'","latidude":'$LONGITUDE',"longitude":'$LONGITUDE'}'
	printf "\n----------------------------------------\n"
	;;

	postSignup)
	printf "\r---------------------------------------\n \
		\rTest for POST /signup -d username':$USERNAME, 'phoneNb':$PHONENB, 'email':$EMAIL, 'password':$PASSWORD\n"
	curl 87.106.79.94:8447/signup -d "{'username':"$USERNAME",'phoneNb':"$PHONENB",'email':"$EMAIL",'password':"$PASSWORD"}"
	printf "\n----------------------------------------\n"
	;;

	*)
	echo "BAD ARGUMENT" $1
	exit 1
	;;
esac
