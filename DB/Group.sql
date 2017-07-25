use forumtemp;

create table _Group (
	GroupName nvarchar(255) not null primary key,
	Untouchable bool,
	Moderator bool,
	ReadOnly bool,
	Admin bool
);
