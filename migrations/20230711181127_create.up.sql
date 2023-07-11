CREATE SCHEMA IF NOT EXISTS EnglishTrainer;

CREATE TABLE EnglishTrainer.Users(  
    id int NOT NULL PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    userName VARCHAR(100) NOT NULL,
    istrainOwnFrases BOOLEAN,
    addedWordsToday int NOT NULL
);

CREATE TABLE EnglishTrainer.EngPhrase(
    id int NOT NULL PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
	addedById int,
    value VARCHAR(150) NOT NULL,
	CONSTRAINT addedById_fk FOREIGN KEY (addedById) REFERENCES EnglishTrainer.Users(id)
);

CREATE TABLE EnglishTrainer.RusPhrase(
    id int NOT NULL PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
	addedById int,
    value VARCHAR(150) NOT NULL,
	CONSTRAINT addedById_fk FOREIGN KEY (addedById) REFERENCES EnglishTrainer.Users(id)
);

CREATE TABLE EnglishTrainer.RusEngPhrase (
	rusId int NOT NULL,
	engId int NOT NULL,
	CONSTRAINT rusId_fk FOREIGN KEY (rusId) REFERENCES EnglishTrainer.RusPhrase(id),
	CONSTRAINT engId_fk FOREIGN KEY (engId) REFERENCES EnglishTrainer.EngPhrase(id)
);

CREATE OR REPLACE FUNCTION EnglishTrainer.AddNewPair(IN rusV VARCHAR(150), IN engV VARCHAR(150), IN addedBy int) RETURNS VARCHAR 
LANGUAGE plpgsql AS
$$
DECLARE
   rus_id int;
   eng_id int;
   old_rus_id int;
   old_eng_id int;
   rus_id_already_exists BOOLEAN;
   eng_id_already_exists BOOLEAN;
BEGIN
	IF NOT EXISTS (SELECT id FROM EnglishTrainer.RusPhrase WHERE value = rusV AND addedById = addedBy) THEN
		INSERT INTO EnglishTrainer.RusPhrase (value, addedById) VALUES (rusV, addedBy) RETURNING id INTO rus_id;
	ELSE
		SELECT id INTO rus_id FROM EnglishTrainer.RusPhrase WHERE value = rusV AND addedById = addedBy;
		rus_id_already_exists = TRUE;
	END IF;
	
	IF NOT EXISTS (SELECT id FROM EnglishTrainer.EngPhrase WHERE value = engV AND addedById = addedBy) THEN
		INSERT INTO EnglishTrainer.EngPhrase (value, addedById) VALUES (engV, addedBy) RETURNING id INTO eng_id;
	ELSE 
		SELECT id INTO eng_id FROM EnglishTrainer.EngPhrase WHERE value = engV AND addedById = addedBy;
		eng_id_already_exists = TRUE;
	END IF;
	IF rus_id_already_exists AND eng_id_already_exists THEN
		IF EXISTS (SELECT rusId FROM EnglishTrainer.RusEngPhrase WHERE rusId = rus_id AND engId = eng_id) THEN
			RETURN 'This pair of phrases already exists';
		END IF;
	END IF;
	INSERT INTO EnglishTrainer.RusEngPhrase (rusId, engId) VALUES (rus_id, eng_id);
	RETURN '';
END
$$;

INSERT INTO EnglishTrainer.Users(userName, istrainOwnFrases, addedWordsToday) VALUES('AdminUser', false, 0);

SELECT EnglishTrainer.AddNewPair ('плохо себя вести', 'to misbehave', 1);
SELECT EnglishTrainer.AddNewPair ('заикаться', 'to stammer', 1);
SELECT EnglishTrainer.AddNewPair ('запоминать', 'to keep in mind', 1);
SELECT EnglishTrainer.AddNewPair ('фразеологизм', 'collocation', 1);
SELECT EnglishTrainer.AddNewPair ('выучить наизусть', 'to learn by heart', 1);
SELECT EnglishTrainer.AddNewPair ('овсянка', 'oatmeal', 1);
SELECT EnglishTrainer.AddNewPair ('быть абсолютно уверенным', 'to be dead sure', 1);
SELECT EnglishTrainer.AddNewPair ('не сыпь мне соль на рану', 'don`t rub salt into the wound', 1);
SELECT EnglishTrainer.AddNewPair ('чайник', 'kettle', 1);
SELECT EnglishTrainer.AddNewPair ('чайник', 'pot', 1);
SELECT EnglishTrainer.AddNewPair ('горшок', 'pot', 1);
SELECT EnglishTrainer.AddNewPair ('домашний мастер', 'handyman', 1);
SELECT EnglishTrainer.AddNewPair ('остатки', 'leftovers', 1);
SELECT EnglishTrainer.AddNewPair ('очень обидно', 'it`s such a shame', 1);
SELECT EnglishTrainer.AddNewPair ('оцеплять', 'to cordon off', 1);
SELECT EnglishTrainer.AddNewPair ('доступным языком', 'In layman`s terms', 1);
SELECT EnglishTrainer.AddNewPair ('языком непрофессионала', 'In layman`s terms', 1);
SELECT EnglishTrainer.AddNewPair ('вездесущий', 'omnipresent', 1);
SELECT EnglishTrainer.AddNewPair ('иметь предрасположенность', 'to have an aptitude to smth', 1);
SELECT EnglishTrainer.AddNewPair ('иметь склонность', 'to have an aptitude to smth', 1);
SELECT EnglishTrainer.AddNewPair ('погружаться во что-либо', 'to delve into smth', 1);
SELECT EnglishTrainer.AddNewPair ('воротник', 'collar', 1);
SELECT EnglishTrainer.AddNewPair ('ошейник', 'collar', 1);
SELECT EnglishTrainer.AddNewPair ('поводок', 'leash', 1);
SELECT EnglishTrainer.AddNewPair ('это по-женски', 'womenlike', 1);
SELECT EnglishTrainer.AddNewPair ('ругать', 'to scold', 1);
SELECT EnglishTrainer.AddNewPair ('быть выписанным', 'to be discharged', 1);
SELECT EnglishTrainer.AddNewPair ('своеобразный', 'peculiar', 1);
SELECT EnglishTrainer.AddNewPair ('самое необходимое', 'bare necessities', 1);
SELECT EnglishTrainer.AddNewPair ('шейный платок', 'handkerchief', 1);
SELECT EnglishTrainer.AddNewPair ('делать вид', 'to make belive', 1);
SELECT EnglishTrainer.AddNewPair ('доносить свою точку зрения', 'to make myself understood', 1);
SELECT EnglishTrainer.AddNewPair ('справляться', 'to make do', 1);
SELECT EnglishTrainer.AddNewPair ('придерживаться диеты', 'to stick to a diet ', 1);
SELECT EnglishTrainer.AddNewPair ('сомневаться', 'to be in two minds', 1);
SELECT EnglishTrainer.AddNewPair ('сомневаться', 'to doubt', 1);
SELECT EnglishTrainer.AddNewPair ('пустоголовый', 'emptyminded', 1);
SELECT EnglishTrainer.AddNewPair ('"классический я"', 'typical of me', 1);
SELECT EnglishTrainer.AddNewPair ('сделайте мне поблажку', 'cut me some slack', 1);
SELECT EnglishTrainer.AddNewPair ('это за счет заведения', 'it`s on the house', 1);
SELECT EnglishTrainer.AddNewPair ('напиться как свинья', 'to get drunk as a skunk', 1);
SELECT EnglishTrainer.AddNewPair ('успешно пройти собеседование', 'to nail a job interview', 1);
SELECT EnglishTrainer.AddNewPair ('ладить, понимать друг друга', 'to see eye to eye', 1);
SELECT EnglishTrainer.AddNewPair ('прилежный', 'diligent', 1);
SELECT EnglishTrainer.AddNewPair ('усидчивый', 'diligent', 1);
SELECT EnglishTrainer.AddNewPair ('язвить', 'to needle', 1);
SELECT EnglishTrainer.AddNewPair ('заранее', 'beforehand', 1);
SELECT EnglishTrainer.AddNewPair ('заранее', 'in advance', 1);
SELECT EnglishTrainer.AddNewPair ('упрямый', 'stubborn', 1);
SELECT EnglishTrainer.AddNewPair ('упрямый', 'stiff-necked', 1);
SELECT EnglishTrainer.AddNewPair ('изгонять', 'expell', 1);
SELECT EnglishTrainer.AddNewPair ('отчислять', 'expell', 1);
SELECT EnglishTrainer.AddNewPair ('обвиняться в убийстве', 'to be charged with murder', 1);
SELECT EnglishTrainer.AddNewPair ('бурчать под нос', 'murmuring under smth breath', 1);
SELECT EnglishTrainer.AddNewPair ('мямля', 'mumbler', 1);
SELECT EnglishTrainer.AddNewPair ('захватывающий', 'breathtaking', 1);
SELECT EnglishTrainer.AddNewPair ('непослушный', 'naughty', 1);
SELECT EnglishTrainer.AddNewPair ('поцелуй меня на прощание', 'kiss me goodbye', 1);
SELECT EnglishTrainer.AddNewPair ('невыносимо', 'unbearable', 1);
SELECT EnglishTrainer.AddNewPair ('задумчивый', 'thoughtful', 1);
SELECT EnglishTrainer.AddNewPair ('подвал', 'cellar', 1);
SELECT EnglishTrainer.AddNewPair ('неотложное дело', 'urgent issue', 1);
SELECT EnglishTrainer.AddNewPair ('развеять сомнения', 'alleviate concerns', 1);
SELECT EnglishTrainer.AddNewPair ('преодолевай себя', 'break your limits', 1);
SELECT EnglishTrainer.AddNewPair ('признаваться', 'make a confession', 1);
SELECT EnglishTrainer.AddNewPair ('в следующий раз мне повезет', 'Better luck next time', 1);
SELECT EnglishTrainer.AddNewPair ('с каждым днем', 'day by day', 1);
SELECT EnglishTrainer.AddNewPair ('свидетель ограбления', 'witness to the robbery', 1);
SELECT EnglishTrainer.AddNewPair ('откровенно', 'frankly', 1);
SELECT EnglishTrainer.AddNewPair ('выбирать наряд', 'choosing an outfit', 1);
SELECT EnglishTrainer.AddNewPair ('запутаться', 'get confused', 1);
SELECT EnglishTrainer.AddNewPair ('обижаться', 'get offended', 1);
SELECT EnglishTrainer.AddNewPair ('быть уволенным', 'get fired', 1);
SELECT EnglishTrainer.AddNewPair ('быть помолвленным', 'get engage', 1);
SELECT EnglishTrainer.AddNewPair ('сгореть на солнце', 'get sunburn', 1);
SELECT EnglishTrainer.AddNewPair ('быть хваленым', 'to be praised', 1);
SELECT EnglishTrainer.AddNewPair ('разозлиться на ', 'get angry with', 1);
SELECT EnglishTrainer.AddNewPair ('Обидно, .. ', 'It is insulting', 1);
SELECT EnglishTrainer.AddNewPair ('И на твоей улице будет праздник', 'your ship will come in to', 1);
SELECT EnglishTrainer.AddNewPair ('Что за шум, что за базар', 'what`s the rumpus', 1);
SELECT EnglishTrainer.AddNewPair ('твоя репутация говорит сама за себя, я "наслышан" о вас', 'your rep precedes you', 1);
SELECT EnglishTrainer.AddNewPair ('то хорошо, то плохо, по-разному', 'strikes and gutters', 1);
SELECT EnglishTrainer.AddNewPair ('оставайся в своем "болоте"', 'stay safe at street level', 1);
SELECT EnglishTrainer.AddNewPair ('тупица', 'shmucks', 1);
SELECT EnglishTrainer.AddNewPair ('потерявши голову по волосам не плачут', 'it`s no use crying over spilled milk', 1);
SELECT EnglishTrainer.AddNewPair ('решительный, упертый человек', 'determened person ', 1);
SELECT EnglishTrainer.AddNewPair ('сходить с ума от страха', 'to go crazy with fear', 1);
SELECT EnglishTrainer.AddNewPair ('спросонья', 'being half-awake', 1);
SELECT EnglishTrainer.AddNewPair ('завораживающий', 'mesmerizing', 1);
SELECT EnglishTrainer.AddNewPair ('оценивать', 'to gauge', 1);
SELECT EnglishTrainer.AddNewPair ('оценивать', 'to estimate', 1);
SELECT EnglishTrainer.AddNewPair ('оценивать', 'to rate', 1);

