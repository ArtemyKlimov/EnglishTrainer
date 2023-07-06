CREATE SCHEMA IF NOT EXISTS EnglishTrainer;

CREATE TABLE EnglishTrainer.Users(  
    id int NOT NULL PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    firstName VARCHAR(100),
    userName VARCHAR(100),
    istrainOwnFrases BOOLEAN,
    addedWordsToday int NOT NULL
);

CREATE TABLE EnglishTrainer.EngPhrase(
    id int NOT NULL PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    value VARCHAR(150) NOT NULL
);

CREATE TABLE EnglishTrainer.RusPhrase(
    id int NOT NULL PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    value VARCHAR(150) NOT NULL
);

CREATE TABLE EnglishTrainer.RusEngPhrase (
	rusId int NOT NULL,
	engId int NOT NULL,
	CONSTRAINT rusId_fk FOREIGN KEY (rusId) REFERENCES EnglishTrainer.RusPhrase(id),
	CONSTRAINT engId_fk FOREIGN KEY (engId) REFERENCES EnglishTrainer.EngPhrase(id)
);

CREATE OR REPLACE FUNCTION EnglishTrainer.AddNewPair(IN rusV VARCHAR(150), IN engV VARCHAR(150)) RETURNS VARCHAR 
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
	IF NOT EXISTS (SELECT id FROM EnglishTrainer.RusPhrase WHERE value = rusV) THEN
		INSERT INTO EnglishTrainer.RusPhrase (value) VALUES (rusV) RETURNING id INTO rus_id;
	ELSE
		SELECT id INTO rus_id FROM EnglishTrainer.RusPhrase WHERE value = rusV;
		rus_id_already_exists = TRUE;
	END IF;
	
	IF NOT EXISTS (SELECT id FROM EnglishTrainer.EngPhrase WHERE value = engV) THEN
		INSERT INTO EnglishTrainer.EngPhrase (value) VALUES (engV) RETURNING id INTO eng_id;
	ELSE 
		SELECT id INTO eng_id FROM EnglishTrainer.EngPhrase WHERE value = engV;
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

SELECT EnglishTrainer.AddNewPair ('плохо себя вести', 'to misbehave');
SELECT EnglishTrainer.AddNewPair ('заикаться', 'to stammer');
SELECT EnglishTrainer.AddNewPair ('запоминать', 'to keep in mind');
SELECT EnglishTrainer.AddNewPair ('фразеологизм', 'collocation');
SELECT EnglishTrainer.AddNewPair ('выучить наизусть', 'to learn by heart');
SELECT EnglishTrainer.AddNewPair ('овсянка', 'oatmeal');
SELECT EnglishTrainer.AddNewPair ('быть абсолютно уверенным', 'to be dead sure');
SELECT EnglishTrainer.AddNewPair ('не сыпь мне соль на рану', 'don`t rub salt into the wound');
SELECT EnglishTrainer.AddNewPair ('чайник', 'kettle');
SELECT EnglishTrainer.AddNewPair ('чайник', 'pot');
SELECT EnglishTrainer.AddNewPair ('горшок', 'pot');
SELECT EnglishTrainer.AddNewPair ('домашний мастер', 'handyman');
SELECT EnglishTrainer.AddNewPair ('остатки', 'leftovers');
SELECT EnglishTrainer.AddNewPair ('очень обидно', 'it`s such a shame');
SELECT EnglishTrainer.AddNewPair ('оцеплять', 'to cordon off');
SELECT EnglishTrainer.AddNewPair ('доступным языком', 'In layman`s terms');
SELECT EnglishTrainer.AddNewPair ('языком непрофессионала', 'In layman`s terms');
SELECT EnglishTrainer.AddNewPair ('вездесущий', 'omnipresent');
SELECT EnglishTrainer.AddNewPair ('иметь предрасположенность', 'to have an aptitude to smth');
SELECT EnglishTrainer.AddNewPair ('иметь склонность', 'to have an aptitude to smth');
SELECT EnglishTrainer.AddNewPair ('погружаться во что-либо', 'to delve into smth');
SELECT EnglishTrainer.AddNewPair ('воротник', 'collar');
SELECT EnglishTrainer.AddNewPair ('ошейник', 'collar');
SELECT EnglishTrainer.AddNewPair ('поводок', 'leash');
SELECT EnglishTrainer.AddNewPair ('это по-женски', 'womenlike');
SELECT EnglishTrainer.AddNewPair ('ругать', 'to scold');
SELECT EnglishTrainer.AddNewPair ('быть выписанным', 'to be discharged');
SELECT EnglishTrainer.AddNewPair ('своеобразный', 'peculiar');
SELECT EnglishTrainer.AddNewPair ('самое необходимое', 'bare necessities');
SELECT EnglishTrainer.AddNewPair ('шейный платок', 'handkerchief');
SELECT EnglishTrainer.AddNewPair ('делать вид', 'to make belive');
SELECT EnglishTrainer.AddNewPair ('доносить свою точку зрения', 'to make myself understood');
SELECT EnglishTrainer.AddNewPair ('справляться', 'to make do');
SELECT EnglishTrainer.AddNewPair ('придерживаться диеты', 'to stick to a diet ');
SELECT EnglishTrainer.AddNewPair ('сомневаться', 'to be in two minds');
SELECT EnglishTrainer.AddNewPair ('сомневаться', 'to doubt');
SELECT EnglishTrainer.AddNewPair ('пустоголовый', 'emptyminded');
SELECT EnglishTrainer.AddNewPair ('"классический я"', 'typical of me');
SELECT EnglishTrainer.AddNewPair ('сделайте мне поблажку', 'cut me some slack');
SELECT EnglishTrainer.AddNewPair ('это за счет заведения', 'it`s on the house');
SELECT EnglishTrainer.AddNewPair ('напиться как свинья', 'to get drunk as a skunk');
SELECT EnglishTrainer.AddNewPair ('успешно пройти собеседование', 'to nail a job interview');
SELECT EnglishTrainer.AddNewPair ('ладить, понимать друг друга', 'to see eye to eye');
SELECT EnglishTrainer.AddNewPair ('прилежный', 'diligent');
SELECT EnglishTrainer.AddNewPair ('усидчивый', 'diligent');
SELECT EnglishTrainer.AddNewPair ('язвить', 'to needle');
SELECT EnglishTrainer.AddNewPair ('заранее', 'beforehand');
SELECT EnglishTrainer.AddNewPair ('заранее', 'in advance');
SELECT EnglishTrainer.AddNewPair ('упрямый', 'stubborn');
SELECT EnglishTrainer.AddNewPair ('упрямый', 'stiff-necked');
SELECT EnglishTrainer.AddNewPair ('изгонять', 'expell');
SELECT EnglishTrainer.AddNewPair ('отчислять', 'expell');
SELECT EnglishTrainer.AddNewPair ('обвиняться в убийстве', 'to be charged with murder');
SELECT EnglishTrainer.AddNewPair ('бурчать под нос', 'murmuring under smth breath');
SELECT EnglishTrainer.AddNewPair ('мямля', 'mumbler');
SELECT EnglishTrainer.AddNewPair ('захватывающий', 'breathtaking');
SELECT EnglishTrainer.AddNewPair ('непослушный', 'naughty');
SELECT EnglishTrainer.AddNewPair ('поцелуй меня на прощание', 'kiss me goodbye');
SELECT EnglishTrainer.AddNewPair ('невыносимо', 'unbearable');
SELECT EnglishTrainer.AddNewPair ('задумчивый', 'thoughtful');
SELECT EnglishTrainer.AddNewPair ('подвал', 'cellar');
SELECT EnglishTrainer.AddNewPair ('неотложное дело', 'urgent issue');
SELECT EnglishTrainer.AddNewPair ('развеять сомнения', 'alleviate concerns');
SELECT EnglishTrainer.AddNewPair ('преодолевай себя', 'break your limits');
SELECT EnglishTrainer.AddNewPair ('признаваться', 'make a confession');
SELECT EnglishTrainer.AddNewPair ('в следующий раз мне повезет', 'Better luck next time');
SELECT EnglishTrainer.AddNewPair ('с каждым днем', 'day by day');
SELECT EnglishTrainer.AddNewPair ('свидетель ограбления', 'witness to the robbery');
SELECT EnglishTrainer.AddNewPair ('откровенно', 'frankly');
SELECT EnglishTrainer.AddNewPair ('выбирать наряд', 'choosing an outfit');
SELECT EnglishTrainer.AddNewPair ('запутаться', 'get confused');
SELECT EnglishTrainer.AddNewPair ('обижаться', 'get offended');
SELECT EnglishTrainer.AddNewPair ('быть уволенным', 'get fired');
SELECT EnglishTrainer.AddNewPair ('быть помолвленным', 'get engage');
SELECT EnglishTrainer.AddNewPair ('сгореть на солнце', 'get sunburn');
SELECT EnglishTrainer.AddNewPair ('быть хваленым', 'to be praised');
SELECT EnglishTrainer.AddNewPair ('разозлиться на ', 'get angry with');
SELECT EnglishTrainer.AddNewPair ('Обидно, .. ', 'It is insulting');
SELECT EnglishTrainer.AddNewPair ('И на твоей улице будет праздник', 'your ship will come in to');
SELECT EnglishTrainer.AddNewPair ('Что за шум, что за базар', 'what`s the rumpus');
SELECT EnglishTrainer.AddNewPair ('твоя репутация говорит сама за себя, я "наслышан" о вас', 'your rep precedes you');
SELECT EnglishTrainer.AddNewPair ('то хорошо, то плохо, по-разному', 'strikes and gutters');
SELECT EnglishTrainer.AddNewPair ('оставайся в своем "болоте"', 'stay safe at street level');
SELECT EnglishTrainer.AddNewPair ('тупица', 'shmucks');
SELECT EnglishTrainer.AddNewPair ('потерявши голову по волосам не плачут', 'it`s no use crying over spilled milk');
SELECT EnglishTrainer.AddNewPair ('решительный, упертый человек', 'determened person ');
SELECT EnglishTrainer.AddNewPair ('сходить с ума от страха', 'to go crazy with fear');
SELECT EnglishTrainer.AddNewPair ('спросонья', 'being half-awake');
SELECT EnglishTrainer.AddNewPair ('завораживающий', 'mesmerizing');
SELECT EnglishTrainer.AddNewPair ('оценивать', 'to gauge');
SELECT EnglishTrainer.AddNewPair ('оценивать', 'to estimate');
SELECT EnglishTrainer.AddNewPair ('оценивать', 'to rate');

