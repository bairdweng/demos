package com.edu.admin.dao;

import java.io.Serializable;
import java.util.Date;

public class FxUser implements Serializable {
    /**
     *
     * This field was generated by MyBatis Generator.
     * This field corresponds to the database column fx_user.id
     *
     * @mbg.generated Fri Nov 15 11:47:09 CST 2019
     */
    private Long id;

    /**
     *
     * This field was generated by MyBatis Generator.
     * This field corresponds to the database column fx_user.user_name
     *
     * @mbg.generated Fri Nov 15 11:47:09 CST 2019
     */
    private String userName;

    /**
     *
     * This field was generated by MyBatis Generator.
     * This field corresponds to the database column fx_user.password
     *
     * @mbg.generated Fri Nov 15 11:47:09 CST 2019
     */
    private String password;

    /**
     *
     * This field was generated by MyBatis Generator.
     * This field corresponds to the database column fx_user.first_name
     *
     * @mbg.generated Fri Nov 15 11:47:09 CST 2019
     */
    private String firstName;

    /**
     *
     * This field was generated by MyBatis Generator.
     * This field corresponds to the database column fx_user.last_name
     *
     * @mbg.generated Fri Nov 15 11:47:09 CST 2019
     */
    private String lastName;

    /**
     *
     * This field was generated by MyBatis Generator.
     * This field corresponds to the database column fx_user.gender
     *
     * @mbg.generated Fri Nov 15 11:47:09 CST 2019
     */
    private String gender;

    /**
     *
     * This field was generated by MyBatis Generator.
     * This field corresponds to the database column fx_user.address
     *
     * @mbg.generated Fri Nov 15 11:47:09 CST 2019
     */
    private String address;

    /**
     *
     * This field was generated by MyBatis Generator.
     * This field corresponds to the database column fx_user.class_name
     *
     * @mbg.generated Fri Nov 15 11:47:09 CST 2019
     */
    private String className;

    /**
     *
     * This field was generated by MyBatis Generator.
     * This field corresponds to the database column fx_user.father_id
     *
     * @mbg.generated Fri Nov 15 11:47:09 CST 2019
     */
    private Long fatherId;

    /**
     *
     * This field was generated by MyBatis Generator.
     * This field corresponds to the database column fx_user.mother_id
     *
     * @mbg.generated Fri Nov 15 11:47:09 CST 2019
     */
    private Long motherId;

    /**
     *
     * This field was generated by MyBatis Generator.
     * This field corresponds to the database column fx_user.creater_time
     *
     * @mbg.generated Fri Nov 15 11:47:09 CST 2019
     */
    private Date createrTime;

    /**
     * This field was generated by MyBatis Generator.
     * This field corresponds to the database table fx_user
     *
     * @mbg.generated Fri Nov 15 11:47:09 CST 2019
     */
    private static final long serialVersionUID = 1L;

    /**
     * This method was generated by MyBatis Generator.
     * This method returns the value of the database column fx_user.id
     *
     * @return the value of fx_user.id
     *
     * @mbg.generated Fri Nov 15 11:47:09 CST 2019
     */
    public Long getId() {
        return id;
    }

    /**
     * This method was generated by MyBatis Generator.
     * This method sets the value of the database column fx_user.id
     *
     * @param id the value for fx_user.id
     *
     * @mbg.generated Fri Nov 15 11:47:09 CST 2019
     */
    public void setId(Long id) {
        this.id = id;
    }

    /**
     * This method was generated by MyBatis Generator.
     * This method returns the value of the database column fx_user.user_name
     *
     * @return the value of fx_user.user_name
     *
     * @mbg.generated Fri Nov 15 11:47:09 CST 2019
     */
    public String getUserName() {
        return userName;
    }

    /**
     * This method was generated by MyBatis Generator.
     * This method sets the value of the database column fx_user.user_name
     *
     * @param userName the value for fx_user.user_name
     *
     * @mbg.generated Fri Nov 15 11:47:09 CST 2019
     */
    public void setUserName(String userName) {
        this.userName = userName;
    }

    /**
     * This method was generated by MyBatis Generator.
     * This method returns the value of the database column fx_user.password
     *
     * @return the value of fx_user.password
     *
     * @mbg.generated Fri Nov 15 11:47:09 CST 2019
     */
    public String getPassword() {
        return password;
    }

    /**
     * This method was generated by MyBatis Generator.
     * This method sets the value of the database column fx_user.password
     *
     * @param password the value for fx_user.password
     *
     * @mbg.generated Fri Nov 15 11:47:09 CST 2019
     */
    public void setPassword(String password) {
        this.password = password;
    }

    /**
     * This method was generated by MyBatis Generator.
     * This method returns the value of the database column fx_user.first_name
     *
     * @return the value of fx_user.first_name
     *
     * @mbg.generated Fri Nov 15 11:47:09 CST 2019
     */
    public String getFirstName() {
        return firstName;
    }

    /**
     * This method was generated by MyBatis Generator.
     * This method sets the value of the database column fx_user.first_name
     *
     * @param firstName the value for fx_user.first_name
     *
     * @mbg.generated Fri Nov 15 11:47:09 CST 2019
     */
    public void setFirstName(String firstName) {
        this.firstName = firstName;
    }

    /**
     * This method was generated by MyBatis Generator.
     * This method returns the value of the database column fx_user.last_name
     *
     * @return the value of fx_user.last_name
     *
     * @mbg.generated Fri Nov 15 11:47:09 CST 2019
     */
    public String getLastName() {
        return lastName;
    }

    /**
     * This method was generated by MyBatis Generator.
     * This method sets the value of the database column fx_user.last_name
     *
     * @param lastName the value for fx_user.last_name
     *
     * @mbg.generated Fri Nov 15 11:47:09 CST 2019
     */
    public void setLastName(String lastName) {
        this.lastName = lastName;
    }

    /**
     * This method was generated by MyBatis Generator.
     * This method returns the value of the database column fx_user.gender
     *
     * @return the value of fx_user.gender
     *
     * @mbg.generated Fri Nov 15 11:47:09 CST 2019
     */
    public String getGender() {
        return gender;
    }

    /**
     * This method was generated by MyBatis Generator.
     * This method sets the value of the database column fx_user.gender
     *
     * @param gender the value for fx_user.gender
     *
     * @mbg.generated Fri Nov 15 11:47:09 CST 2019
     */
    public void setGender(String gender) {
        this.gender = gender;
    }

    /**
     * This method was generated by MyBatis Generator.
     * This method returns the value of the database column fx_user.address
     *
     * @return the value of fx_user.address
     *
     * @mbg.generated Fri Nov 15 11:47:09 CST 2019
     */
    public String getAddress() {
        return address;
    }

    /**
     * This method was generated by MyBatis Generator.
     * This method sets the value of the database column fx_user.address
     *
     * @param address the value for fx_user.address
     *
     * @mbg.generated Fri Nov 15 11:47:09 CST 2019
     */
    public void setAddress(String address) {
        this.address = address;
    }

    /**
     * This method was generated by MyBatis Generator.
     * This method returns the value of the database column fx_user.class_name
     *
     * @return the value of fx_user.class_name
     *
     * @mbg.generated Fri Nov 15 11:47:09 CST 2019
     */
    public String getClassName() {
        return className;
    }

    /**
     * This method was generated by MyBatis Generator.
     * This method sets the value of the database column fx_user.class_name
     *
     * @param className the value for fx_user.class_name
     *
     * @mbg.generated Fri Nov 15 11:47:09 CST 2019
     */
    public void setClassName(String className) {
        this.className = className;
    }

    /**
     * This method was generated by MyBatis Generator.
     * This method returns the value of the database column fx_user.father_id
     *
     * @return the value of fx_user.father_id
     *
     * @mbg.generated Fri Nov 15 11:47:09 CST 2019
     */
    public Long getFatherId() {
        return fatherId;
    }

    /**
     * This method was generated by MyBatis Generator.
     * This method sets the value of the database column fx_user.father_id
     *
     * @param fatherId the value for fx_user.father_id
     *
     * @mbg.generated Fri Nov 15 11:47:09 CST 2019
     */
    public void setFatherId(Long fatherId) {
        this.fatherId = fatherId;
    }

    /**
     * This method was generated by MyBatis Generator.
     * This method returns the value of the database column fx_user.mother_id
     *
     * @return the value of fx_user.mother_id
     *
     * @mbg.generated Fri Nov 15 11:47:09 CST 2019
     */
    public Long getMotherId() {
        return motherId;
    }

    /**
     * This method was generated by MyBatis Generator.
     * This method sets the value of the database column fx_user.mother_id
     *
     * @param motherId the value for fx_user.mother_id
     *
     * @mbg.generated Fri Nov 15 11:47:09 CST 2019
     */
    public void setMotherId(Long motherId) {
        this.motherId = motherId;
    }

    /**
     * This method was generated by MyBatis Generator.
     * This method returns the value of the database column fx_user.creater_time
     *
     * @return the value of fx_user.creater_time
     *
     * @mbg.generated Fri Nov 15 11:47:09 CST 2019
     */
    public Date getCreaterTime() {
        return createrTime;
    }

    /**
     * This method was generated by MyBatis Generator.
     * This method sets the value of the database column fx_user.creater_time
     *
     * @param createrTime the value for fx_user.creater_time
     *
     * @mbg.generated Fri Nov 15 11:47:09 CST 2019
     */
    public void setCreaterTime(Date createrTime) {
        this.createrTime = createrTime;
    }

    /**
     * This method was generated by MyBatis Generator.
     * This method corresponds to the database table fx_user
     *
     * @mbg.generated Fri Nov 15 11:47:09 CST 2019
     */
    @Override
    public String toString() {
        StringBuilder sb = new StringBuilder();
        sb.append(getClass().getSimpleName());
        sb.append(" [");
        sb.append("Hash = ").append(hashCode());
        sb.append(", id=").append(id);
        sb.append(", userName=").append(userName);
        sb.append(", password=").append(password);
        sb.append(", firstName=").append(firstName);
        sb.append(", lastName=").append(lastName);
        sb.append(", gender=").append(gender);
        sb.append(", address=").append(address);
        sb.append(", className=").append(className);
        sb.append(", fatherId=").append(fatherId);
        sb.append(", motherId=").append(motherId);
        sb.append(", createrTime=").append(createrTime);
        sb.append(", serialVersionUID=").append(serialVersionUID);
        sb.append("]");
        return sb.toString();
    }
}