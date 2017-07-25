--use <YOUR DATABASE NAME>;

drop table if exists _Message;
drop table if exists _Theme;
drop table if exists _Subforum;
drop table if exists _User;
drop table if exists _Group;
drop table if exists _Forum;

create table _Forum (
	Id int AUTO_INCREMENT not null,
	Name nvarchar(255) not null,
	primary key(Id)
);

create table _Group (
	Name nvarchar(255) not null primary key,
	Untouchable bool,
	Moderator bool,
	ReadOnly bool,
	Admin bool
);

create table _User (
	Nickname nvarchar(255) primary key,
	Info nvarchar(255),
	Hash binary(32) not null,
	Cookie nvarchar(255),
	GroupName nvarchar(255) not null,
	foreign key (GroupName) references _Group(Name)
);

create table _Subforum (
	Id int not null auto_increment primary key,
	Name nvarchar(255),
	Description nvarchar(255),
	ForumId int not null,
	foreign key (ForumId) references _Forum(Id)
);

create table _Theme (
	Id int not null auto_increment primary key,
	Name nvarchar(255) not null,
	Description nvarchar(255),
	Creator nvarchar(255) not null,
	SubforumId int not null,
	foreign key (Creator) references _User(Nickname),
	foreign key (SubforumId) references _Subforum(Id)
);

create table _Message (
	Id int not null auto_increment primary key,
	_Text nvarchar(255) not null,
	_Date datetime not null,
	Creator nvarchar(255) not null,
	ThemeId int not null,
	foreign key (Creator) references _User(Nickname),
	foreign key (ThemeId) references _Theme(Id)
);

insert into _Group values ('Admin', true,true,false,true);
insert into _Group values ('Moderator', false, true, false, false);
insert into _Group values ('DELETED', true, false, true, false);
insert into _Group values ('DEFAULT', false, false, false, false);
insert into _User values ('DELETED', '',MD5('dfgdjfgdg73__'),'','DELETED');
--insert into _User values ('<YOUR ADMIN NICKNAME>','<YOUR ADMIN INFO>',MD5('<YOUR ADMIN PASSWORD>'),'','Admin');
