<?php

class MemcachedApi
{
	protected $_memcache;

	public function __construct()
	{
		$this->_memcache = new \Memcache();
		$this->_memcache->connect('127.0.0.1', 3000);
	}

	public function getStats()
	{
		return $this->_memcache->getStats();
	}

	public function getUserById($userId)
	{
		return $this->_memcache->get($this->_command("GetUserById", (int) $userId));
	}

	public function getUserByTwoParams($login, $userId)
	{
		return $this->_memcache->get($this->_command("GetUserByTwoParams", $login, (int) $userId));
	}

	public function getAuthUser($token)
	{
		return $this->_memcache->get($this->_command("GetAuthUser", $token));
	}

	public function cast()
	{
		return $this->_memcache->get($this->_command("Cast", 42, 3.14159265359, "Hello"));
	}

	public function multiGet()
	{
		return $this->_memcache->get(
			[
				$this->_command("GetUserById", 42),
				$this->_command("GetAuthUser", "token")
			]
		);
	}
	
	public function returnError()
	{
		return $this->_memcache->get($this->_command("ReturnError"));
	}

	public function notFoundmethod()
	{
		return $this->_memcache->get($this->_command("notFoundmethod"));
	}

	public function GetUserWhereIdIn()
	{
		return $this->_memcache->get($this->_command("GetUserWhereIdIn", [1,2,3,4]));
	}

	public function setUser() 
	{
		return $this->_memcache->set(
			'SetUser', json_encode(
				[
					'user_id'    => 42, 
					'user_name'  => 'New User', 
					'user_login' => 'new_login'
				]
			)
		);
	}

	public function increment($delta) 
	{
		return $this->_memcache->increment("Increment", $delta);
	}

	public function decrement($delta) 
	{
		return $this->_memcache->decrement("Decrement", $delta);
	}

	public function delete()
	{
		return $this->_memcache->delete($this->_command("Delete", 42));
	}

	protected function _command($method, ...$params)
	{
		return sprintf("%s:%s", $method, base64_encode(json_encode($params)));
	}
}

$Api = new MemcachedApi;

var_dump($Api->getStats());
var_dump($Api->getUserById(42));

var_dump($Api->multiGet());

var_dump($Api->notFoundmethod());
var_dump($Api->getAuthUser(uniqid()));
var_dump($Api->getUserByTwoParams(uniqid(), 42));
var_dump($Api->cast());
var_dump($Api->returnError());
var_dump($Api->notFoundmethod());
var_dump($Api->GetUserWhereIdIn());

var_dump($Api->setUser());

for ($i = 1; $i <= 5; $i++) { 
	
	var_dump($Api->increment($i));
}

for ($i = 1; $i <= 5; $i++) { 
	
	var_dump($Api->decrement($i));
}

var_dump($Api->delete());