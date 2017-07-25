use forumtemp;

create table _Subforum (
	SubforumID int not null auto_increment primary key,
	Name nvarchar(255),
	Description nvarchar(255),
	_Number int not null,
	F_ID int not null,
	foreign key (F_ID) references _Forum(ForumID)
);
