part of 'chat_bloc.dart';

abstract class ChatEvent extends Equatable {
  const ChatEvent();

  @override
  List<Object> get props => [];
}

class LoadChatHistory extends ChatEvent {
  final String matchId;
  const LoadChatHistory({required this.matchId});
  @override
  List<Object> get props => [matchId];
}

class ConnectWebSocket extends ChatEvent {}

class DisconnectWebSocket extends ChatEvent {}

class SendMessage extends ChatEvent {
  final String matchId;
  final String targetId;
  final String content;
  final String type;

  const SendMessage({
    required this.matchId,
    required this.targetId,
    required this.content,
    this.type = 'text',
  });

  @override
  List<Object> get props => [matchId, targetId, content, type];
}

class ReceiveMessage extends ChatEvent {
  final ChatMessage message;
  const ReceiveMessage({required this.message});
  @override
  List<Object> get props => [message];
}
