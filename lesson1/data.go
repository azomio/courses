package main

type Movie struct {
	ID     int
	Name   string
	Poster string
}

func MovieList() []*Movie {
	return []*Movie{
		&Movie{0, "Бойцовски клуб", "/static/posters/fightclub.jpg"},
		&Movie{1, "Крестный отец", "/static/posters/father.jpg"},
		&Movie{2, "Железный человек", "/static/posters/ironman.jpg"},
		&Movie{3, "Холодное сердце", "/static/posters/cold.jpg"},
		&Movie{4, "Эверест", "/static/posters/everest.jpg"},
		&Movie{5, "Криминальное чтиво", "/static/posters/pulpfiction.jpg"},
		&Movie{6, "Пираты карибского моря", "/static/posters/pirates.jpg"},
		&Movie{7, "Оно", "/static/posters/it.jpg"},
		&Movie{8, "Однажды в голливуде", "/static/posters/hw.jpg"},
		&Movie{9, "Звездные войны", "/static/posters/starwars.jpg"},
	}
}

type User struct {
	ID       int
	Login    string
	Password string
	Session  string
}

func UserList() []*User {
	return []*User{
		&User{0, "bob", "god", "21298df8a3277357ee55b01df9530b535cf08ec1"},
		&User{1, "alice", "secret", "e5e9fa1ba31ecd1ae84f75caaa474f3a663f05f4"},
	}
}

func LoadUserBySession(ses string) *User {
	uu := UserList()

	for _, u := range uu {
		if ses == u.Session {
			return u
		}
	}

	return nil
}

func LoadUserByLogin(login string) *User {
	uu := UserList()

	for _, u := range uu {
		if login == u.Login {
			return u
		}
	}

	return nil
}
