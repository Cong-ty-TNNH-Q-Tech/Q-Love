import 'dart:convert';
import 'package:http/http.dart' as http;
import '../domain/entities/wingman_suggestion.dart';

abstract class ChatRemoteDataSource {
  Future<List<WingmanSuggestion>> getWingmanSuggestions(String matchId);
}

class ChatRemoteDataSourceImpl implements ChatRemoteDataSource {
  final http.Client client;

  ChatRemoteDataSourceImpl({required this.client});

  @override
  Future<List<WingmanSuggestion>> getWingmanSuggestions(String matchId) async {
    // In a real app, baseUrl is injected via env
    final url = Uri.parse('http://localhost:3000/api/v1/matches/$matchId/wingman-suggestions');
    
    final response = await client.get(
      url,
      headers: {'Content-Type': 'application/json'},
    );

    if (response.statusCode == 200) {
      final jsonResponse = json.decode(response.body);
      final List<dynamic> suggestionsJson = jsonResponse['suggestions'] ?? [];
      return suggestionsJson.map((json) => WingmanSuggestion.fromJson(json)).toList();
    } else {
      throw Exception('Failed to load wingman suggestions');
    }
  }
}
