# Gerador-autounattend

**Gerador de arquivos `autounattend.xml` para instalações automatizadas do Windows 10 e Windows 11.**

O **Gerador-autounattend** é um aplicativo para Windows que permite criar arquivos `autounattend.xml` de maneira simples e intuitiva, por meio de uma interface gráfica com mais de 400 opções de personalização.

Com ele, você configura previamente diversos aspectos da instalação e do sistema operacional, eliminando a necessidade de interagir manualmente com as telas de configuração do Windows.

O projeto utiliza como referência um modelo baseado no estilo do [Winhance](https://github.com/memstechtips/Winhance), oferecendo ampla liberdade de personalização antes da geração do arquivo XML.

![Interface do Gerador-autounattend](https://github.com/user-attachments/assets/bb7feeaa-eac6-4a49-b8dc-0b9d94ab7bd3)

## ✨ Principais recursos

- **Mais de 400 opções configuráveis:** personalização dividida em diferentes categorias e abas.
- **Presets e perfis personalizados:** utilize configurações predefinidas ou salve seus próprios perfis para reutilizá-los posteriormente.
- **Interface gráfica intuitiva:** configure todas as opções sem precisar editar arquivos XML manualmente.
- **Visualização prévia:** consulte o conteúdo do arquivo XML antes de salvá-lo.
- **Carregamento de configurações:** importe perfis anteriormente salvos.
- **Execução local:** não exige conexão com a internet para gerar os arquivos.

## ⚙️ Funcionalidades

### Personalização do Windows

Configure diferentes aspectos do sistema operacional, incluindo:

- Sistema e configurações gerais.
- Privacidade e segurança.
- Desempenho e otimizações.
- Personalização da interface.
- Remoção de aplicativos pré-instalados.
- Outras configurações avançadas do Windows.

### Contas e instalação

- Criação e configuração de múltiplas contas de usuário.
- Definição de contas como Administrador ou Usuário comum.
- Configuração de logon automático.
- Personalização do papel de parede por URL ou script PowerShell.
- Definição de idioma, região e fuso horário.
- Configurações adicionais do processo de instalação.

### Pós-instalação

O aplicativo permite gerar automaticamente um script para instalação de programas após a conclusão da instalação do Windows.

Recursos disponíveis:

- Instalação automatizada de programas utilizando **WinGet** e **Chocolatey**.
- Suporte a listas personalizadas de aplicativos.
- Opção para instalar automaticamente o gerenciador Chocolatey.
- Criação do script diretamente na Área de Trabalho.

### Segurança e prevenção de riscos

Para reduzir a possibilidade de configurações indesejadas:

- Opções potencialmente arriscadas são identificadas visualmente.
- Alertas são apresentados antes da geração do arquivo.
- A opção de apagamento do disco exige uma confirmação adicional.
- Foram removidas configurações que poderiam comprometer a conexão com a internet ou o funcionamento do Wi-Fi.
- A opção de instalação do Winhance e seu respectivo atalho na Área de Trabalho foram removidos.

## 🎨 Interface e personalização

O aplicativo conta com uma interface moderna, leve e organizada, seguindo uma identidade visual em azul e branco.

**Características:**

- Modo claro e modo escuro.
- Três idiomas disponíveis:
  - Português (Brasil).
  - Inglês.
  - Espanhol.
- Preferências de idioma e tema salvas automaticamente.
- Ícone próprio do aplicativo.
- Interface adaptada para facilitar a configuração das opções.

## 🚀 Como funciona

O Gerador-autounattend funciona de maneira independente, sem necessidade de instalação.

1. Execute o arquivo `.exe`.
2. O aplicativo abrirá sua interface gráfica em uma janela própria.
3. Configure as opções desejadas nas diferentes abas.
4. Se preferir, selecione um preset ou carregue um perfil salvo anteriormente.
5. Visualize o conteúdo do XML antes de finalizar.
6. Gere e salve o arquivo `autounattend.xml` utilizando a janela nativa do Windows.

**Tecnologia de execução:** o aplicativo utiliza o Microsoft Edge ou o Google Chrome, quando disponíveis, para apresentar sua interface. Caso nenhum deles seja encontrado, utiliza o navegador padrão do sistema.

Todo o processamento de geração do XML é realizado localmente, sem necessidade de conexão com a internet.

## 💻 Compatibilidade

| Requisito | Compatibilidade |
|---|---|
| Windows 10 | Todas as versões |
| Windows 11 | Todas as versões, incluindo 26H2 |
| Arquitetura x64 | Suportada |
| Arquitetura x86 (32 bits) | Suportada |
| Arquitetura ARM64 | Suportada |

O arquivo XML gerado é destinado a instalações do Windows que oferecem suporte ao mecanismo de instalação automatizada por meio do `autounattend.xml`.

## 📌 Observações importantes

- O arquivo gerado deve ser utilizado em uma mídia de instalação compatível do Windows.
- Revise cuidadosamente as configurações antes de iniciar uma instalação automatizada.
- Opções que envolvem exclusão de disco ou alterações críticas devem ser utilizadas com atenção.
- A disponibilidade e o comportamento de determinadas configurações podem variar conforme a versão e a edição do Windows.

## 🙏 Créditos

Projeto desenvolvido com base em um arquivo de configuração gerado pela ferramenta [Winhance](https://github.com/memstechtips/Winhance), criada por **memstechtips**.

A ferramenta foi adaptada e expandida com uma interface própria, novas possibilidades de configuração e recursos adicionais para facilitar a criação de instalações automatizadas do Windows.

---

**Gerador-autounattend** — Simplifique, personalize e automatize a instalação do Windows.