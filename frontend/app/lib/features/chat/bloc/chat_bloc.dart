// Copyright (c) 2026 Q-Tech. All rights reserved.
// Licensed under the GNU AGPLv3 License.

import 'dart:async';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:equatable/equatable.dart';
import '../data/chat_model.dart';
import '../data/chat_repository.dart';

part 'chat_event.dart';
part 'chat_state.dart';

class ChatBloc extends Bloc<ChatEvent, ChatState> {
  final ChatRepository repository;
  StreamSubscription? _messageSubscription;

  ChatBloc({required this.repository}) : super(ChatInitial()) {
    on<LoadChatHistory>(_onLoadChatHistory);
    on<SendMessage>(_onSendMessage);
    on<ReceiveMessage>(_onReceiveMessage);
    on<ConnectWebSocket>(_onConnectWebSocket);
    on<DisconnectWebSocket>(_onDisconnectWebSocket);
  }

  Future<void> _onLoadChatHistory(LoadChatHistory event, Emitter<ChatState> emit) async {
    emit(ChatLoading());
    try {
      final history = await repository.getHistory(event.matchId);
      // Generate some mock history if empty for demo purposes
      emit(ChatLoaded(messages: history));
    } catch (e) {
      emit(ChatError(message: e.toString()));
    }
  }

  void _onConnectWebSocket(ConnectWebSocket event, Emitter<ChatState> emit) {
    repository.connectWebSocket();
    _messageSubscription = repository.messageStream.listen(
      (message) {
        add(ReceiveMessage(message: message));
      },
      onError: (error) {
        emit(ChatError(message: 'WebSocket Error: $error'));
      },
    );
  }

  void _onReceiveMessage(ReceiveMessage event, Emitter<ChatState> emit) {
    if (state is ChatLoaded) {
      final currentState = state as ChatLoaded;
      final updatedMessages = List<ChatMessage>.from(currentState.messages)..insert(0, event.message);
      emit(ChatLoaded(messages: updatedMessages));
    }
  }

  void _onSendMessage(SendMessage event, Emitter<ChatState> emit) {
    if (state is ChatLoaded) {
      final message = ChatMessage(
        id: DateTime.now().millisecondsSinceEpoch.toString(),
        matchId: event.matchId,
        senderId: repository.currentUserId,
        targetId: event.targetId,
        type: event.type,
        content: event.content,
        createdAt: DateTime.now(),
        isMine: true,
      );
      
      // Optimistic UI update
      final currentState = state as ChatLoaded;
      final updatedMessages = List<ChatMessage>.from(currentState.messages)..insert(0, message);
      emit(ChatLoaded(messages: updatedMessages));

      // Send to server
      try {
        repository.sendMessage(message);
      } catch (e) {
        // Rollback optimistic update
        final rollbackMessages = List<ChatMessage>.from(updatedMessages)..remove(message);
        emit(ChatLoaded(messages: rollbackMessages));
        emit(ChatError(message: 'Failed to send message: $e'));
      }
    }
  }

  void _onDisconnectWebSocket(DisconnectWebSocket event, Emitter<ChatState> emit) {
    _messageSubscription?.cancel();
    repository.disconnect();
  }

  @override
  Future<void> close() {
    _messageSubscription?.cancel();
    repository.disconnect();
    return super.close();
  }
}
