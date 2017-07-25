use forumtemp;

create table _Theme (
	ThemeID int not null auto_increment primary key,
	Name nvarchar(255) not null,
	Description nvarchar(255),
	_Number int not null,
	FMessage nvarchar(255) not null,
	Creator nvarchar(255) not null,
	S_ID int not null,
	foreign key (Creator) references _User(Nickname),
	foreign key (S_ID) references _Subforum(SubforumID)
);
