// credential store HTTP handlers
// credentials in form
//  {
//      "type": <credential type>
//      "name": <credential name>
//      "data": {
//          <credential data>
//      }
//  }
// credentials can be fetched by searched by name, filtered by type
// data portion is encrypted
//
package credentialstore

import (
    "github.com/gin-gonic/gin"
)

// Controller method to add new credential
func addCredential(c *gin.Context) {
    
}

// Controller method to list all credentials
func listCredentials(c *gin.Context) {
    
}

// Controller method to fetch specific credential
func getCredential(c *gin.Context) {
    
}

// Controller method to remove specific credential
func deleteCredential(c *gin.Context) {
    
}

// 
func RunHttpContollers() {
    router := gin.Default()
    
    v1 := router.Group("/credential/creds")
    {
        v1.POST("/", addCredential)
        v1.GET("/", listCredentials)
        v1.GET("/:id", getCredential)
        v1.DELETE("/:id", deleteCredential)
    }
    
    router.Run()
}