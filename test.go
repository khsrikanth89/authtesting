package main
import (    
     "fmt"
	 "net/http"
	"appengine"
    "crypto/sha256"
    "appengine/datastore"
	"io/ioutil"
    "encoding/json"
    "time"
)
func init(){
   http.HandleFunc("/signup", signup)
   http.HandleFunc("/login",login)
}
    
type Data struct{
	Id string
	Password string
    Email string
}

type Session struct{
    UserId string
    Id string
    // Device string
    // DeviceId string
    Expiry time.Time    
}

func Hash256(input string) (output string) {
    hash := sha256.New()
    bytepass := []byte(input)
    hash.Write(bytepass)
    sum := hash.Sum(nil)
    return fmt.Sprintf("%x", sum)
}

func login(response http.ResponseWriter, request *http.Request){
     context := appengine.NewContext(request)
      data := make(map[string]string)
         if request.Method == "OPTIONS" {
        response.Header().Set("Access-Control-Allow-Methods","GET, HEAD, PUT, DELETE, POST")
        response.Header().Set("Access-Control-Allow-Headers","origin, accept, content-type")
        response.Header().Set("Access-Control-Allow-Origin","http://127.0.0.1:8081")
        response.Header().Set("Access-Control-Allow-Credentials","true")
        response.Header().Add("X-Requested-With", "XMLHttpRequest")
        response.WriteHeader(200)
        return
    }else if request.Method == "POST"{
            response.Header().Set("Access-Control-Allow-Origin","http://127.0.0.1:8081")
            response.Header().Set("Access-Control-Allow-Credentials","true")
            var f interface{}
            jsonBinInfo, _ := ioutil.ReadAll(request.Body)
            request.Body.Close()
            json.Unmarshal(jsonBinInfo, &f)
            incomingData := f.(map[string]interface{})
            for i, v := range incomingData {
                if v != "" {
                    data[i] = v.(string)
                }
            }
    }
        queryResult := make([]Data, 0, 100)
        q := datastore.NewQuery("Users").Filter("Id =", data["id"]).Filter("Password =",Hash256(data["password"]))
        key,_ := q.GetAll(context, &queryResult)
        context.Infof(fmt.Sprint("%v",request))
        if len(queryResult)==0{
            response.WriteHeader(422)
            response.Write([]byte(`{"error":"Invalid Username/Password"}`))
        }else{
            session := http.Cookie{
                Name:"session",
                Value: data["id"] + fmt.Sprint(key[0].IntID),
                Path: "/",
                Expires:  time.Now().Add(356 * 24 * time.Hour),
            }


        s1 := Session {
                UserId: fmt.Sprint(key[0].IntID),
                Id: data["id"],
                Expiry: time.Now().Add(356 * 24 * time.Hour),
        }
        datastore.Put(context, datastore.NewKey(context, "Session", "", 0, nil), &s1)
             http.SetCookie(response,&session)
             response.WriteHeader(200)
             response.Write([]byte(`{"success":"Logged In"}`))
        }
}

func signup(response http.ResponseWriter, request *http.Request){
   context := appengine.NewContext(request)
    data := make(map[string]string)
    if request.Method == "OPTIONS" {
        response.Header().Set("Access-Control-Allow-Methods","GET, HEAD, PUT, DELETE, POST")
        response.Header().Set("Access-Control-Allow-Headers","origin, accept, content-type")
        response.Header().Set("Access-Control-Allow-Origin","http://127.0.0.1:8081")
        response.Header().Set("Access-Control-Allow-Credentials","true")
        response.Header().Add("X-Requested-With", "XMLHttpRequest")
        response.WriteHeader(200)
        return
    }else if request.Method == "POST"{
            response.Header().Set("Access-Control-Allow-Origin","http://127.0.0.1:8081")
            response.Header().Set("Access-Control-Allow-Credentials","true")
            var f interface{}
            jsonBinInfo, _ := ioutil.ReadAll(request.Body)
            request.Body.Close()
            json.Unmarshal(jsonBinInfo, &f)
            incomingData := f.(map[string]interface{})
            for i, v := range incomingData {
                if v != "" {
                    data[i] = v.(string)
                }
            }
    }
    var userExist bool
    queryResult := make([]Data, 0, 100)
    if data["id"]!=""{
        q := datastore.NewQuery("Users").Filter("Id =", data["id"])
        q.GetAll(context, &queryResult)
        if len(queryResult)==0{
            userExist = false
        }else{
            userExist = true
        }
    } 
    if data["email"]!=""{
        context.Infof("Inside Email");
        q := datastore.NewQuery("Users").Filter("Email =", data["email"])
        q.GetAll(context, &queryResult)
        context.Infof(fmt.Sprintf("%+v",q))
        if len(queryResult)==0{
            userExist = false
        }else{
            userExist = true
        }
    }
    if userExist == false{
        d1 := Data {
            Id: data["id"],
            Password: Hash256(data["password"]), 
            Email: data["email"],
        }

        context.Infof(fmt.Sprintf("%+v",d1))
        _, err := datastore.Put(context, datastore.NewKey(context, "Users", "", 0, nil), &d1)
        if err != nil {
            http.Error(response, err.Error(), http.StatusInternalServerError)
            return
        }
    }else{
        response.WriteHeader(422)
        response.Write([]byte(`{"error":"Username/Emailid Taken"}`))
    }
}