package main
import (    
     "fmt"
	 "net/http"
	"appengine"
    "crypto/sha256"
    "appengine/datastore"
    "golang.org/x/oauth2/google"
	"io/ioutil"
    "encoding/json"
    "appengine/urlfetch"
    "net/url"
    "time"
    "appengine/blobstore"
    "appengine/image"
    cloudstore "google.golang.org/cloud/storage"
)
func init(){
   http.HandleFunc("/signup", signup)
   http.HandleFunc("/login",login)
   http.HandleFunc("/profile",profile)
   http.HandleFunc("/avatar",avatar)
   http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){ http.ServeFile(w,r,"test.html")})
}
    
type Data struct{
	Id string
	Password string
    Email string
}

type Avatar struct{
    Url string
    Thumbnail string
}   

type Profile struct{
    Id string
    Name string
    Age string
    Phone string
    Avatars []Avatar 
}


type Session struct{
    UserId string
    Id string
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
        // response.Header().Set("Access-Control-Allow-Origin","http://127.0.0.1:8081")
        response.Header().Set("Access-Control-Allow-Credentials","true")
        response.Header().Add("X-Requested-With", "XMLHttpRequest")
        response.WriteHeader(200)
        return
    }else if request.Method == "POST"{
            // response.Header().Set("Access-Control-Allow-Origin","http://127.0.0.1:8081")
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
                Value: data["id"],
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
        // response.Header().Set("Access-Control-Allow-Origin","http://127.0.0.1:8081")
        response.Header().Set("Access-Control-Allow-Credentials","true")
        response.Header().Add("X-Requested-With", "XMLHttpRequest")
        response.WriteHeader(200)
        return
    }else if request.Method == "POST"{
            // response.Header().Set("Access-Control-Allow-Origin","http://127.0.0.1:8081")
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

func profile(response http.ResponseWriter, request *http.Request){
    context := appengine.NewContext(request)
    data := make(map[string]string)
        if request.Method == "OPTIONS" {
            response.Header().Set("Access-Control-Allow-Methods","GET, HEAD, PUT, DELETE, POST")
            response.Header().Set("Access-Control-Allow-Headers","origin, accept, content-type")
            // response.Header().Set("Access-Control-Allow-Origin","http://127.0.0.1:8081")
            response.Header().Set("Access-Control-Allow-Credentials","true")
            response.Header().Add("X-Requested-With", "XMLHttpRequest")
            response.WriteHeader(200)
            return
        }else if request.Method == "POST"{
                // response.Header().Set("Access-Control-Allow-Origin","http://127.0.0.1:8081")
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
                context.Infof(fmt.Sprintf("%+v",data))
                queryResult := make([]Data, 0, 100)
                session, err := request.Cookie("session")
                if err!=nil{
                    return
                }
                context.Infof(fmt.Sprintf("%+v",session))
                q := datastore.NewQuery("Session").Filter("Id =", session.Value)
                _,_ = q.GetAll(context, &queryResult)
                if len(queryResult)==0{
                    response.WriteHeader(422)
                    response.Write([]byte(`{"error":"No Session"}`))
                    return
                }
        
                prof := Profile{
                Id: session.Value,
                Name:data["name"],
                Age: data["age"],
                Phone: data["phone"],
                }

                _, err = datastore.Put(context, datastore.NewKey(context, "Profile", "", 0, nil), &prof)
                if err != nil {
                    http.Error(response, err.Error(), http.StatusInternalServerError)
                    return
                }
          }else if request.Method == "GET"{
                response.Header().Set("Access-Control-Allow-Credentials","true")
                qr := make([]Profile, 0, 100)
                session, _ := request.Cookie("session")
                q := datastore.NewQuery("Profile").Filter("Id =", session.Value)
                _,_ = q.GetAll(context, &qr)
                response.Write([]byte(fmt.Sprintf("%+v",qr[0])))
                return
            }      
}

func avatar(response http.ResponseWriter, request *http.Request){
    context := appengine.NewContext(request)
    queryResult := make([]Data, 0, 100)
    _,fileHead,_ := request.FormFile("content")
    // data := make(map[string]string)
    if request.Method == "OPTIONS" {
        response.Header().Set("Access-Control-Allow-Methods","GET, HEAD, PUT, DELETE, POST")
        response.Header().Set("Access-Control-Allow-Headers","origin, accept, content-type")
        // response.Header().Set("Access-Control-Allow-Origin","http://127.0.0.1:8081")
        response.Header().Set("Access-Control-Allow-Credentials","true")
        response.Header().Add("X-Requested-With", "XMLHttpRequest")
        response.WriteHeader(200)
            return
    }else if request.Method == "POST"{
        // response.Header().Set("Access-Control-Allow-Origin","http://127.0.0.1:8081")
        response.Header().Set("Access-Control-Allow-Credentials","true")
        // var f interface{}
        session, err := request.Cookie("session")
        if err!=nil{
            return
        }
        context.Infof(fmt.Sprintf("%+v",session))
        q := datastore.NewQuery("Session").Filter("Id =", session.Value)
        _,_ = q.GetAll(context, &queryResult)
        if len(queryResult)==0{
            response.WriteHeader(422)
            response.Write([]byte(`{"error":"No Session"}`))
            return
        }
        avatar,_ := sendData(session.Value, request)
        thumbnail,_ := makeAvatar(context, "/gs/test-auth-service/"+session.Value+"/"+fileHead.Filename)
        qresult := make([]Profile,0,100)
        context.Infof(fmt.Sprintf("%+v",avatar))
        q = datastore.NewQuery("Profile").Filter("Id =", session.Value)
        key,_ := q.GetAll(context, &qresult)
        context.Infof(fmt.Sprintf("%+v",qresult))
            ava := Avatar{
                Url:avatar,
                Thumbnail:thumbnail,
            }
            var temp []Avatar
            temp = append(temp,ava)
        if len(qresult)==0{
            prof := Profile{
                Id: session.Value,
                Avatars:temp,
            }
         _, err = datastore.Put(context, datastore.NewKey(context, "Profile", "", 0, nil), &prof)
        }else{
             qresult[0].Avatars = append(qresult[0].Avatars,ava)
             _, err = datastore.Put(context, key[0], &qresult[0])
        }
        if err != nil {
            http.Error(response, err.Error(), http.StatusInternalServerError)
            return
        }
        context.Infof(fmt.Sprintf("%+v",qresult[0]))
        response.Write([]byte(avatar))
    }
}

func sendData(user string, r *http.Request) (string, error) {
    file, fileHeader, ferr := r.FormFile("content")
    c := appengine.NewContext(r)
    if ferr != nil {
        c.Infof("ferr :%v", ferr)
        return "", ferr
    }
    client := urlfetch.Client(c)
    c.Infof("Request :%v", r)
    var requrl url.URL
    token,err := google.AppEngineTokenSource(c,
    cloudstore.ScopeFullControl).Token()

    requrl.Scheme = "https"
    requrl.Host = "www.googleapis.com"
    requrl.Path = "/upload/storage/v1/b/test-auth-service/o"
    param := url.Values{}
    param.Set("uploadType", "media")
    param.Set("name", user+"/"+fileHeader.Filename)
    requrl.RawQuery = param.Encode()
    reqst, _ := http.NewRequest("POST", requrl.String(), file)
    reqst.Header.Set("Content-Type", fileHeader.Header["Content-Type"][0])
    reqst.Header.Set("Authorization", "Bearer "+ token.AccessToken)
    reqst.Header.Set("Content-Length", r.Header.Get("Content-Length"))
    c.Infof("Outgoing :%v", reqst)
    resp, err := client.Do(reqst)
    c.Infof(fmt.Sprintf("%v",resp))
    if err != nil {
        c.Infof("ERR : ", err)
        return "", err
    } else {
        avatar, _ := url.Parse("http://storage.googleapis.com/test-auth-service/"+user+"/"+fileHeader.Filename)
        return avatar.String(), nil
    }
}

func makeAvatar(context appengine.Context, path string) (string, error) {
    key, _ := blobstore.BlobKeyForFile(context,path)
    opts := image.ServingURLOptions{Secure:false, Size:50, Crop:false}
    url, _ := image.ServingURL(context, key, &opts)
    context.Infof("url : ", url.String())
    return url.String(), nil
}