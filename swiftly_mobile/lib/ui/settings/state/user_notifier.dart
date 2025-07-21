import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:swiftly_mobile/domain/models/label_item.dart';

import '../../../domain/user/models/user.dart';
import '../../core/themes/colors.dart';
import 'user_state.dart';

class UserNotifier extends StateNotifier<UserState> {
  UserNotifier() : super(UserState.initial());

  void addUser(User user) {
    state = state.copyWith(users: [...state.users, user]);
  }

  void removeUser(String id) {
    state = state.copyWith(
      users: state.users.where((user) => user.id != id).toList(),
    );
  }

  void updateName(String userId, String newName) {
    state = state.copyWith(
      users:
          state.users.map((user) {
            return user.id == userId ? user.copyWith(name: newName) : user;
          }).toList(),
    );
  }

  void updateLastName(String userId, String newLastName) {
    state = state.copyWith(
      users:
          state.users.map((user) {
            return user.id == userId ? user.copyWith(name: newLastName) : user;
          }).toList(),
    );
  }

  void updateImage(String userId, String newImage) {
    state = state.copyWith(
      users:
          state.users.map((user) {
            return user.id == userId ? user.copyWith(name: newImage) : user;
          }).toList(),
    );
  }

  void updateRole(String userId, String newRole) {
    state = state.copyWith(
      users:
          state.users.map((user) {
            return user.id == userId ? user.copyWith(name: newRole) : user;
          }).toList(),
    );
  }

  void loadUsers() async {
    state = state.copyWith(isLoading: true);
    await Future.delayed(const Duration(seconds: 1));
    state = state.copyWith(
      users: [...state.users, ...mockUsers],
      isLoading: false,
    );
  }
}

final mockUsers = [
  User.create(id: 'hhh', name: 'Павел'),
  User.create(
    id: 'bbb',
    name: 'Вася',
    lastName: 'Пупкин',
    role: LabelItem(title: 'программист', color: AppColors.amaranthMagenta),
  ),
  User.create(
    id: 'ccc',
    name: 'Алеша',
    lastName: 'Попович',
    role: LabelItem(title: 'дизайнер', color: AppColors.wanderingThrus),
  ),
  User.create(
    id: 'ddd',
    name: 'Анастасия',
    lastName: 'Петровна',
    role: LabelItem(title: 'аналитик', color: AppColors.wanderingThrus),
  ),
  User.create(
    id: 'eee',
    name: 'Алексей',
    lastName: 'Семеныч',
    role: LabelItem(title: 'программист', color: AppColors.amaranthMagenta),
  ),
];
