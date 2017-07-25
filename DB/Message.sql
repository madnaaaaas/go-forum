use forumtemp;

create table _Message (
	MessageID int not null auto_increment primary key,
	_Text nvarchar(255) not null,
	_Date datetime not null,
	_Number int not null,
	Creator nvarchar(255) not null,
	T_ID int not null,
	foreign key (Creator) references _User(Nickname),
	foreign key (T_ID) references _Theme(ThemeID)
);

