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
		return $this->_memcache->get($this->_getCommand("GetUserById", (int) $userId));
	}

	public function getUserByTwoParams($login, $userId)
	{
		return $this->_memcache->get($this->_getCommand("GetUserByTwoParams", $login, (int) $userId));
	}

	public function getAuthUser($token)
	{
		return $this->_memcache->get($this->_getCommand("GetAuthUser", $token));
	}

	public function cast()
	{
		return $this->_memcache->get($this->_getCommand("Cast", 42, 3.14159265359, "Hello"));
	}

	public function returnError()
	{
		return $this->_memcache->get($this->_getCommand("ReturnError"));
	}

	public function notFoundmethod()
	{
		return $this->_memcache->get($this->_getCommand("notFoundmethod"));
	}

	protected function _getCommand($method, ...$params)
	{
		return sprintf("%s:%s", $method, base64_encode(json_encode($params)));
	}
}

$Api = new MemcachedApi;

var_dump($Api->getStats());
var_dump($Api->getUserById(42));
var_dump($Api->notFoundmethod());
var_dump($Api->getAuthUser(uniqid()));
var_dump($Api->getUserByTwoParams(uniqid(), 42));
var_dump($Api->cast());
var_dump($Api->returnError());
var_dump($Api->notFoundmethod());


