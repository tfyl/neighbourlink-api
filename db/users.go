package db

import (
	"neighbourlink-api/types"
)

func (db *DB) AddUser (user types.User) (types.User, error) {

	tx := db.MustBegin()
	_ ,err  := tx.Exec(`
		WITH ins1 AS (
		   INSERT INTO user_detail(username, email)
		   VALUES ($1, $2)
		   RETURNING user_id
		   )
		, ins2 AS (
		   INSERT INTO user_auth (user_id, password,permissions)
		   SELECT user_id, $3,$4 FROM ins1
		   RETURNING user_id
		   )
		INSERT INTO user_attribute (user_id, local_area,reputation)
		SELECT user_id, $5, $6 FROM ins2;
	`,user.Username,user.Email,user.Password,user.Permissions,user.LocalArea,user.Reputation)
	if err != nil {
		return user, err
	}

	err = tx.Commit()
	return user, err
}



func (db *DB) GetUserAll () ([]types.User, error) {
	var users []types.User
	err := db.Select(&users, `
	SELECT
	       user_detail.user_id,
	       user_detail.username,
	       user_detail.email,
	       user_auth.password,
	       user_auth.permissions,
	       user_attribute.local_area,
	       user_attribute.reputation
	       
	FROM
	     user_detail
	         
	INNER JOIN user_auth ON user_detail.user_id=user_auth.user_id
	INNER JOIN user_attribute ON user_detail.user_id=user_attribute.user_id
	;`, )
	if err != nil {
		return users,err
	}

	return users,err

}


func (db *DB) GetUserByID (user types.User) (types.User, error) {


	err := db.Get(&user, `
	SELECT
	       user_detail.user_id,
	       user_detail.username,
	       user_detail.email,
	       user_auth.password,
	       user_auth.permissions,
	       user_attribute.local_area,
	       user_attribute.reputation
	       
	FROM
	     user_detail
	         
	INNER JOIN user_auth ON user_detail.user_id=user_auth.user_id
	INNER JOIN user_attribute ON user_detail.user_id=user_attribute.user_id
	
	WHERE
		user_detail.user_id=$1;


	;`, user.UserID)
	if err != nil {
		return user,err
	}

	return user,err

}

func (db *DB) GetUserByUsername (user types.User) (types.User, error) {

	err := db.Get(&user, `
	SELECT
	       user_detail.user_id,
	       user_detail.username,
	       user_detail.email,
	       user_auth.password,
	       user_auth.permissions,
	       user_attribute.local_area,
	       user_attribute.reputation
	       
	FROM
	     user_detail
	         
	INNER JOIN user_auth ON user_detail.user_id=user_auth.user_id
	INNER JOIN user_attribute ON user_detail.user_id=user_attribute.user_id
	
	WHERE
		user_detail.username=$1;


	;`, user.Username)
	if err != nil {
		return user,err
	}

	return user,err

}

func (db *DB) UpdateUser (user types.User) (types.User, error) {



	tx := db.MustBegin()
	// update user_detail table
	_ ,err  := tx.Exec(`
		UPDATE 
		        user_detail 
		SET 
		        username = $2,
		        email = $3
		WHERE
				user_id = $1
				`, user.UserID,user.Username,user.Email)
	if err != nil {
		return user,err
	}
	// update user_auth table

	_ ,err  = tx.Exec(`
		UPDATE 
		        user_auth 
		SET 
		        password = $2,
		        permissions = $3
		WHERE
				user_id = $1
				`, user.UserID,user.Password,user.Permissions)
	if err != nil {
		return user,err
	}
	// update user_attribute table
	_ ,err  = tx.Exec(`
		UPDATE 
		        user_attribute 
		SET 
		        local_area = $2,
		        reputation = $3
		WHERE
				user_id = $1
				`, user.UserID,user.LocalArea,user.Reputation)
	if err != nil {
		return user,err
	}

	err = tx.Commit()

	return user,err

}


func (db *DB) DeleteUser (user types.User) (types.User, error) {

	tx := db.MustBegin()
	// update user_detail table
	_ ,err  := tx.Exec(`
		DELETE FROM 
		        user_detail 
		WHERE
				user_id = $1
				`, user.UserID)
	if err != nil {
		return user,err
	}
	// update user_auth table

	_ ,err  = tx.Exec(`
		DELETE FROM  
		        user_auth 
		WHERE
				user_id = $1
				`, user.UserID)
	if err != nil {
		return user,err
	}
	// update user_attribute table
	_ ,err  = tx.Exec(`
		DELETE FROM  
		        user_attribute 
		WHERE
				user_id = $1
				`, user.UserID)
	if err != nil {
		return user,err
	}

	err = tx.Commit()

	return user,err

}