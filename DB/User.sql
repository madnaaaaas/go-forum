use forumtemp;

create table _User (
	Nickname nvarchar(255) primary key,
	User_info nvarchar(255),
	Gr_ID nvarchar(255) not null,
	foreign key (Gr_ID) references _Group(GroupName)
);
	
