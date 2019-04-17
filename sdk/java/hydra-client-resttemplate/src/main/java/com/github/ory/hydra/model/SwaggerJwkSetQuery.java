/*
 * ORY Hydra
 * Welcome to the ORY Hydra HTTP API documentation. You will find documentation for all HTTP APIs here.
 *
 * OpenAPI spec version: latest
 * 
 *
 * NOTE: This class is auto generated by the swagger code generator program.
 * https://github.com/swagger-api/swagger-codegen.git
 * Do not edit the class manually.
 */


package com.github.ory.hydra.model;

import java.util.Objects;
import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonValue;
import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;

/**
 * SwaggerJwkSetQuery
 */
@javax.annotation.Generated(value = "io.swagger.codegen.languages.JavaClientCodegen", date = "2019-04-17T17:33:18.188+02:00")
public class SwaggerJwkSetQuery {
  @JsonProperty("set")
  private String set = null;

  public SwaggerJwkSetQuery set(String set) {
    this.set = set;
    return this;
  }

   /**
   * The set in: path
   * @return set
  **/
  @ApiModelProperty(required = true, value = "The set in: path")
  public String getSet() {
    return set;
  }

  public void setSet(String set) {
    this.set = set;
  }


  @Override
  public boolean equals(java.lang.Object o) {
    if (this == o) {
      return true;
    }
    if (o == null || getClass() != o.getClass()) {
      return false;
    }
    SwaggerJwkSetQuery swaggerJwkSetQuery = (SwaggerJwkSetQuery) o;
    return Objects.equals(this.set, swaggerJwkSetQuery.set);
  }

  @Override
  public int hashCode() {
    return Objects.hash(set);
  }


  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class SwaggerJwkSetQuery {\n");
    
    sb.append("    set: ").append(toIndentedString(set)).append("\n");
    sb.append("}");
    return sb.toString();
  }

  /**
   * Convert the given object to string with each line indented by 4 spaces
   * (except the first line).
   */
  private String toIndentedString(java.lang.Object o) {
    if (o == null) {
      return "null";
    }
    return o.toString().replace("\n", "\n    ");
  }
  
}

