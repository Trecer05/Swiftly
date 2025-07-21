import 'package:freezed_annotation/freezed_annotation.dart';

import '../../../domain/user/models/user.dart';

part 'user_state.freezed.dart';

@freezed
class UserState with _$UserState {
  const factory UserState({
    required List<User> users,
    @Default(false) bool isLoading,
  }) = _UserState;

  factory UserState.initial() => const UserState(users: []);
}
