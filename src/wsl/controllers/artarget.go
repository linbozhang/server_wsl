package controllers

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"wsl/models"
	"log"
	"time"
	"encoding/base64"
	"encoding/hex"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"bytes"
	"net/http"
	"github.com/astaxie/beego"
)
type VuforiaImageData struct{
	Name string `json:"name"`
	Width float32	`json:"width"`
	Image string	`json:"image"`
	Active_flag int	`json:"active_flag"`
	Application_metadata string `json:"application_metadata"`
}
// ArtargetController operations for Artarget
type ArtargetController struct {
	beego.Controller
}
const publickey ="b558528870f94d19585d8e41cc23adaf711c7ca6"
const secretkey="5967f67529553d6f483e8f55bda544a6757d7200"
const empty="d41d8cd98f00b204e9800998ecf8427e"
type VuforiaController struct {
	beego.Controller
}

func getMd5(content []byte ) string{
	h:=md5.New()
	h.Write(content)
	cipherStr:=h.Sum(nil)
	return hex.EncodeToString( cipherStr)
}

func Sign(content string) string {
	//h:=sha1.New()
	//io.WriteString(h,content)
	key:=[]byte(secretkey)
	mac:=hmac.New(sha1.New,key)
	mac.Write([]byte(content))
	 return base64.StdEncoding.EncodeToString( mac.Sum(nil))
}
func (c *VuforiaController) URLMapping(){
	c.Mapping("Post", c.Post)
}
func (c *VuforiaController) Post(){
	// brandname:=c.GetString("brandname")
	// log.Printf("brandname:",brandname)
	// f,h,err:=c.GetFile("uploadname")
	
	// if err!=nil{
	// 	log.Printf("rawbrand getfile err",err)
	// }
	// defer f.Close();
	// path:="ar/public/target/"+"target_"+h.Filename
	// c.SaveToFile("uploadname",path)
	//  v:=new(models.Brand)
	//  v.Id=0
	//  v.LogoUrl="http://localhost:8080/"+path
	//  v.Name=brandname
	//  if _, err := models.AddBrand(v); err == nil {
	// 	c.Ctx.Output.SetStatus(201)
	// 	c.Data["json"] = v
	// } else {
	// 	c.Data["json"] = err.Error()
	// }
	// c.ServeJSON()


	
	f,h,err:=c.GetFile("uploadname")
		
	if err!=nil{
		log.Printf("getfile err",err)
	}
	defer f.Close();

	log.Printf(""+h.Filename)
	path:="ar/public/target/"+"target_"+h.Filename
	c.SaveToFile("uploadname",path)
	dat,err:=ioutil.ReadFile(path)
	if(err!=nil){
		log.Printf("read file error",err)
	}

	tempName:="wsl_"+strconv.FormatInt( time.Now().Unix(),10)
	log.Printf("tempname",tempName)
	targetdata:=VuforiaImageData{Name:tempName,Width:1,Image:base64.StdEncoding.EncodeToString(dat),Active_flag:0,Application_metadata:base64.StdEncoding.EncodeToString([] byte (tempName))}

	body,err:=json.Marshal(targetdata)
	if(err !=nil){
		log.Printf("json err:",err)
	}
	log.Printf(string(body))
	
	//tm := time.Unix(time.Now().Unix(), 0)
	datetime:= time.Now().UTC().Format(time.RFC1123)  //tm.Format("2006-01-02 03:04:05 PM")
	// datetime
	datetime=strings.Replace(datetime,"UTC","GMT",1)
	log.Printf(datetime)
	md5result:=getMd5(body)
	stringToSign:="POST"+"\n"+md5result+"\n"+"application/json"+"\n"+datetime+"\n"+"/targets"
	signresult:=Sign(stringToSign)
	url:="https://vws.vuforia.com/targets"
	req,err:=http.NewRequest("POST",url,bytes.NewBuffer(body))
	req.Header.Set("Content-Type","application/json")
	req.Header.Set("Authorization","VWS "+ publickey +":"+signresult)
	req.Header.Set("Date", datetime)
	client:=&http.Client{}
	resp,err:=client.Do(req)
		
	if err!=nil{
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println("response Status:", resp.Status)
    fmt.Println("response Headers:", resp.Header)
    returnBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(returnBody))

	c.Data["json"]=string(returnBody)
	c.ServeJSON()
}
// URLMapping ...
func (c *ArtargetController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Post
// @Description create Artarget
// @Param	body		body 	models.Artarget	true		"body for Artarget content"
// @Success 201 {int} models.Artarget
// @Failure 403 body is empty
// @router / [post]
func (c *ArtargetController) Post() {
	var v models.Artarget
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if _, err := models.AddArtarget(&v); err == nil {
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = v
		} else {
			c.Data["json"] = err.Error()
		}
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// GetOne ...
// @Title Get One
// @Description get Artarget by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Artarget
// @Failure 403 :id is empty
// @router /:id [get]
func (c *ArtargetController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetArtargetById(id)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}

// GetAll ...
// @Title Get All
// @Description get Artarget
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Artarget
// @Failure 403
// @router / [get]
func (c *ArtargetController) GetAll() {
	var fields []string
	var sortby []string
	var order []string
	var query = make(map[string]string)
	var limit int64 = 10
	var offset int64

	// fields: col1,col2,entity.col3
	if v := c.GetString("fields"); v != "" {
		fields = strings.Split(v, ",")
	}
	// limit: 10 (default is 10)
	if v, err := c.GetInt64("limit"); err == nil {
		limit = v
	}
	// offset: 0 (default is 0)
	if v, err := c.GetInt64("offset"); err == nil {
		offset = v
	}
	// sortby: col1,col2
	if v := c.GetString("sortby"); v != "" {
		sortby = strings.Split(v, ",")
	}
	// order: desc,asc
	if v := c.GetString("order"); v != "" {
		order = strings.Split(v, ",")
	}
	// query: k:v,k:v
	if v := c.GetString("query"); v != "" {
		for _, cond := range strings.Split(v, ",") {
			kv := strings.SplitN(cond, ":", 2)
			if len(kv) != 2 {
				c.Data["json"] = errors.New("Error: invalid query key/value pair")
				c.ServeJSON()
				return
			}
			k, v := kv[0], kv[1]
			query[k] = v
		}
	}

	l, err := models.GetAllArtarget(query, fields, sortby, order, offset, limit)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = l
	}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the Artarget
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Artarget	true		"body for Artarget content"
// @Success 200 {object} models.Artarget
// @Failure 403 :id is not int
// @router /:id [put]
func (c *ArtargetController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.Artarget{Id: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if err := models.UpdateArtargetById(&v); err == nil {
			c.Data["json"] = "OK"
		} else {
			c.Data["json"] = err.Error()
		}
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the Artarget
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *ArtargetController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	if err := models.DeleteArtarget(id); err == nil {
		c.Data["json"] = "OK"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}
