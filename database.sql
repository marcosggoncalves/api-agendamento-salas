DROP DATABASE IF EXISTS agendamentos;

CREATE DATABASE agendamentos;

use agendamentos;

CREATE TABLE usuario (
  id INTEGER UNSIGNED  NOT NULL   AUTO_INCREMENT,
  nome VARCHAR(255)  NULL,
  email VARCHAR(50)  NULL,
  senha VARCHAR(255)  NULL,
  cpf VARCHAR(50)  NULL,
  created_at DATETIME  NULL,
  updated_at DATETIME  NULL,
  deleted_at DATETIME  NULL,
PRIMARY KEY(id));

CREATE TABLE sala (
  id INTEGER UNSIGNED  NOT NULL   AUTO_INCREMENT,
  hash VARCHAR(255)  NULL,
  nome VARCHAR(255)  NULL,
  descricao TEXT  NULL,
  intervalo_por_agendamento INT  NULL,
  horario_ini_funcionamento TIME  NULL,
  horario_fim_funcionamento TIME  NULL,
  created_at DATETIME  NULL,
  updated_at DATETIME  NULL,
  deleted_at DATETIME  NULL,
PRIMARY KEY(id));

CREATE TABLE departamento (
  id INTEGER UNSIGNED  NOT NULL   AUTO_INCREMENT,
  descricao TEXT  NULL,
  created_at DATETIME  NULL,
  updated_at DATETIME  NULL,
  deleted_at DATETIME  NULL,
PRIMARY KEY(id));

CREATE TABLE sala_grade_horario (
  id INTEGER UNSIGNED  NOT NULL   AUTO_INCREMENT,
  sala_id INTEGER UNSIGNED  NOT NULL,
  inicial TIME  NULL,
  final TIME  NULL, 
  created_at DATETIME  NULL,
  updated_at DATETIME  NULL,
  deleted_at DATETIME  NULL,
PRIMARY KEY(id),
  FOREIGN KEY(sala_id)
    REFERENCES sala(id)
      ON DELETE CASCADE
      ON UPDATE CASCADE );

CREATE TABLE sala_grade_horario_reserva (
  id INTEGER UNSIGNED  NOT NULL   AUTO_INCREMENT,
  sala_grade_horario_id INTEGER UNSIGNED  NOT NULL,
  departamento_id INTEGER UNSIGNED  NOT NULL,
  nome VARCHAR(255)  NULL, 
    status VARCHAR(100)  NULL DEFAULT "RESERVADO",
  data_reserva DATE  NULL,
  created_at DATETIME  NULL,
  updated_at DATETIME  NULL,
  deleted_at DATETIME  NULL,
PRIMARY KEY(id),
  FOREIGN KEY(departamento_id)
    REFERENCES departamento(id)
      ON DELETE CASCADE
      ON UPDATE CASCADE,
  FOREIGN KEY(sala_grade_horario_id)
    REFERENCES sala_grade_horario(id)
      ON DELETE CASCADE
      ON UPDATE CASCADE);

INSERT INTO `usuario` (`id`, `nome`, `cpf`, `email`, `senha`, `created_at`, `updated_at`, `deleted_at`) VALUES
(1, 'Marcos Lopes', '069.389.071-10', 'marcos.ggoncalves.pr@gmail.com', '$2a$10$E/x3r7OGpuDkrPj/KtHyM.YceWWJ4fi7ME2etxY/9IdfMYDo9T9dW', '2024-11-10 19:54:34', '2024-11-10 19:54:34', NULL);

INSERT INTO `sala` (`id`, `hash`, `nome`, `descricao`, `intervalo_por_agendamento`, `horario_ini_funcionamento`, `horario_fim_funcionamento`, `created_at`, `updated_at`, `deleted_at`) VALUES
(1, 'uQ1hlDYLigFw4bPu', 'Sala 01 - Sala Cooperativa', 'A maior sala de reunião com mesa e cadeiras.', 30, '08:00:00', '18:00:00', '2024-11-13 14:57:46', '2024-11-13 22:33:50', NULL),
(2, 'ov1Wu0fh4Ebzp52E', 'Sala 02 - Sala NEO', 'Sala com sofás e TV interativa.', 30, '08:00:00', '18:00:00', '2024-11-13 20:16:44', '2024-11-13 22:24:26', NULL),
(3, 'xBBvhv36C_0Re5Kp', 'Sala 03 - Sala Verde', 'Sala ao lado da diretoria, com sofás e TV.', 30, '08:00:00', '18:00:00', '2024-11-13 20:17:08', '2024-11-13 20:46:53', NULL),
(4, 'DUB5m4ii0e9Uq0AF', 'Sala 04 - Sala Maringá', 'Sala menor com mesa e cadeiras ao lado do comercial.', 30, '08:00:00', '18:00:00', '2024-11-13 20:17:28', '2024-11-13 20:46:58', NULL),
(5, '4QnqortQdJpw6wLn', 'Carro Comercial - Onix', 'Carro usado pela equipe comercial para fechamento comercial', 30, '08:00:00', '18:00:00', '2024-11-13 20:17:52', '2024-11-13 22:24:10', NULL);

INSERT INTO `departamento` (`id`, `descricao`, `created_at`, `updated_at`, `deleted_at`) VALUES
(1, 'Tecnologia', '2024-11-13 11:17:32', '2024-11-13 22:34:39', NULL),
(2, 'RH', '2024-11-13 11:17:37', '2024-11-13 11:17:37', NULL),
(3, 'DP', '2024-11-13 11:17:43', '2024-11-13 11:17:43', NULL),
(4, 'Diretoria', '2024-11-13 11:17:48', '2024-11-13 11:17:48', NULL),
(5, 'Cooperativa', '2024-11-13 11:18:02', '2024-11-13 11:18:02', NULL),
(6, 'Compras', '2024-11-13 11:18:08', '2024-11-13 11:18:08', NULL),
(7, 'Financeiro', '2024-11-13 11:18:15', '2024-11-13 11:18:15', NULL),
(8, 'Comercial', '2024-11-13 11:18:21', '2024-11-13 11:18:21', NULL),
(9, 'Marketing', '2024-11-13 11:18:28', '2024-11-13 11:18:28', NULL),
(10, 'Outros', '2024-11-13 11:18:41', '2024-11-13 11:18:41', NULL),
(11, 'Engenharia', '2024-11-13 11:18:46', '2024-11-13 11:18:46', NULL),
(12, 'Recepção', '2024-11-13 20:41:31', '2024-11-13 20:41:31', NULL),
(13, 'Administrativo', '2024-11-13 20:41:43', '2024-11-13 22:35:08', NULL);

ALTER TABLE sala ADD COLUMN color VARCHAR(255);

UPDATE sala SET color = '#FF0000' WHERE id = 1;  -- Vermelho
UPDATE sala SET color = '#0000FF' WHERE id = 2;  -- Azul
UPDATE sala SET color = '#008000' WHERE id = 3;  -- Verde
UPDATE sala SET color = '#FFFF00' WHERE id = 4;  -- Amarelo
UPDATE sala SET color = '#FFA500' WHERE id = 5;  -- Laranja 

ALTER TABLE sala_grade_horario ADD COLUMN dia_inteiro int DEFAULT 0;

INSERT INTO `sala_grade_horario`(`sala_id`, `inicial`, `final`, `created_at`, `updated_at`, `dia_inteiro`) 
SELECT s.id as sala_id, s.horario_ini_funcionamento as inicial, s.horario_fim_funcionamento as final, now() as created_at, now() as updated_at, 1 as dia_inteiro FROM sala s


delete  FROM `sala_grade_horario_reserva` WHERE id in (228,229,230,231,232,233,234,235,236,237,238,239,246,243,245,240,247,244,242,241)